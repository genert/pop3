package pop3

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net/textproto"
)

// Client holds the net conn and read/write buffer objects.
type Connection struct {
	conn   io.ReadWriteCloser
	Reader *textproto.Reader
	Writer *textproto.Writer
}

// NewConnection initializes a connection.
func NewConnection(conn io.ReadWriteCloser) *Connection {
	return &Connection{
		conn,
		textproto.NewReader(bufio.NewReader(conn)),
		textproto.NewWriter(bufio.NewWriter(conn)),
	}
}

// Close closes a connection.
func (c *Connection) Close() error {
	return c.conn.Close()
}

// Cmd sends the given command on the connection.
func (c *Connection) Cmd(format string, args ...interface{}) (string, error) {
	if err := c.Writer.PrintfLine(format, args...); err != nil {
		return "", fmt.Errorf("failed to write with format and args: %w", err)
	}

	return c.ReadLine()
}

// ReadLine reads a single line from the buffer.
func (c *Connection) ReadLine() (string, error) {
	line, err := c.Reader.ReadLine()
	if err != nil {
		return "", fmt.Errorf("failed to read line: %w", err)
	}

	if len(line) < 1 {
		return "", errors.New("empty response")
	}

	if IsErr(line) {
		return line, fmt.Errorf("something went wrong: %s", line)
	}

	return line, nil
}

// ReadLines reads from the buffer until it hits the message end dot (".").
func (c *Connection) ReadLines() (lines []string, err error) {
	for {
		line, err := c.ReadLine()
		if err != nil {
			return nil, fmt.Errorf("failed to read line: %w", err)
		}

		// Look for a dot to indicate the end of a message
		// from the server.
		if line == "." {
			break
		}
		lines = append(lines, line)
	}
	return
}
