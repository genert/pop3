package pop3

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/jhillyerd/enmime"
)

// MessageList represents the metadata returned by the server for a
// message stored in the maildrop.
type MessageList struct {
	// Non unique id reported by the server
	ID int

	// Size of the message
	Size int
}

const (
	CommandReset = "RSET"

	// CommandStat is a command to retrieve statistics about mailbox.
	CommandStat = "STAT"

	// CommandDelete is a command to delete message from POP3 server.
	CommandDelete = "DELE"

	// CommandList is a command to get list of messages from POP3 server.
	CommandList = "LIST"

	// CommandNoop is a ping-like command that tells POP3 to do nothing.
	// (i.e. send something line pong-response).
	CommandNoop = "NOOP"

	// CommandPassword is a command to send user password to POP3 server.
	CommandPassword = "PASS"

	// CommandQuit is a command to tell POP3 server that you are quitting.
	CommandQuit = "QUIT"

	// CommandRetrieve is a command to retrieve POP3 message from server.
	CommandRetrieve = "RETR"

	// CommandUser is a command to send user login to POP3 server.
	CommandUser = "USER"
)

// Client for POP3.
type Client struct {
	conn *Connection
}

// Dial opens new connection and creates a new POP3 client.
func Dial(addr string) (c *Client, err error) {
	var conn net.Conn
	if conn, err = net.Dial("tcp", addr); err != nil {
		return nil, fmt.Errorf("failed to dial: %w", err)
	}

	return NewClient(conn)
}

// DialTLS opens new TLS connection and creates a new POP3 client.
func DialTLS(addr string) (c *Client, err error) {
	var conn *tls.Conn
	if conn, err = tls.Dial("tcp", addr, nil); err != nil {
		return nil, fmt.Errorf("failed to dial tls: %w", err)
	}
	return NewClient(conn)
}

// NewClient creates a new POP3 client.
func NewClient(conn net.Conn) (*Client, error) {
	c := &Client{
		conn: NewConnection(conn),
	}

	// Make sure we receive the server greeting
	line, err := c.conn.ReadLine()
	if err != nil {
		return nil, fmt.Errorf("failed to read line: %w", err)
	}

	if !IsOK(line) {
		return nil, fmt.Errorf("server did not response with +OK: %s", line)
	}

	return c, nil
}

// Authorization logs into POP3 server with login and password.
func (c *Client) Authorization(user, pass string) error {
	if _, err := c.conn.Cmd("%s %s", CommandUser, user); err != nil {
		return fmt.Errorf("failed at USER command: %w", err)
	}

	if _, err := c.conn.Cmd("%s %s", CommandPassword, pass); err != nil {
		return fmt.Errorf("failed at PASS command: %w", err)
	}

	return c.Noop()
}

// Quit sends the QUIT message to the POP3 server and closes the connection.
func (c *Client) Quit() error {
	if _, err := c.conn.Cmd(CommandQuit); err != nil {
		return fmt.Errorf("failed at QUIT command: %w", err)
	}

	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %w", err)
	}

	return nil
}

// Noop will do nothing however can prolong the end of a connection.
func (c *Client) Noop() error {
	if _, err := c.conn.Cmd(CommandNoop); err != nil {
		return fmt.Errorf("failed at NOOP command: %w", err)
	}

	return nil
}

// Stat retrieves a drop listing for the current maildrop, consisting of the
// number of messages and the total size (in octets) of the maildrop.
// In the event of an error, all returned numeric values will be 0.
func (c *Client) Stat() (count, size int, err error) {
	line, err := c.conn.Cmd(CommandStat)
	if err != nil {
		return
	}

	if len(strings.Fields(line)) != 3 {
		return 0, 0, fmt.Errorf("invalid response returned from server: %s", line)
	}

	// Number of messages in maildrop
	count, err = strconv.Atoi(strings.Fields(line)[1])
	if err != nil {
		return
	}
	if count == 0 {
		return
	}

	// Total size of messages in bytes
	size, err = strconv.Atoi(strings.Fields(line)[2])
	if err != nil {
		return
	}
	if size == 0 {
		return
	}
	return
}

// ListAll returns a MessageList object which contains all messages in the maildrop.
func (c *Client) ListAll() (list []MessageList, err error) {
	if _, err = c.conn.Cmd(CommandList); err != nil {
		return
	}

	lines, err := c.conn.ReadLines()
	if err != nil {
		return
	}

	for _, v := range lines {
		id, err := strconv.Atoi(strings.Fields(v)[0])
		if err != nil {
			return nil, err
		}

		size, err := strconv.Atoi(strings.Fields(v)[1])
		if err != nil {
			return nil, err
		}
		list = append(list, MessageList{id, size})
	}
	return
}

// Rset will unmark any messages that have being marked for deletion in
// the current session.
func (c *Client) Rset() error {
	if _, err := c.conn.Cmd(CommandReset); err != nil {
		return fmt.Errorf("failed at RSET command: %w", err)
	}
	return nil
}

// Retr downloads the given message and returns it as a mail.Message object.
func (c *Client) Retr(msg int) (*enmime.Envelope, error) {
	if _, err := c.conn.Cmd("%s 1", CommandRetrieve); err != nil {
		return nil, fmt.Errorf("failed at RETR command: %w", err)
	}

	message, err := enmime.ReadEnvelope(c.conn.Reader.DotReader())
	if err != nil {
		return nil, fmt.Errorf("failed to read message: %w", err)
	}

	return message, nil
}

// Dele will delete the given message from the maildrop.
// Changes will only take affect after the Quit command is issued.
func (c *Client) Dele(msg int) error {
	if _, err := c.conn.Cmd("%s %d", CommandDelete, msg); err != nil {
		return fmt.Errorf("failed at DELE command: %w", err)
	}
	return nil
}
