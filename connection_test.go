package pop3_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/genert/pop3"
)

func TestConnection_Cmd(t *testing.T) {
	testCases := []struct {
		response string
		command  string
	}{
		{
			"+OK 13 messages:\r\n",
			"STAT",
		},
		{
			"+OK\r\n",
			"LIST 1",
		},
	}

	for _, tt := range testCases {
		var buffer []byte
		buffer = make([]byte, 0)
		var fake readWriteFaker
		fake.Reader = strings.NewReader(tt.response)
		fake.Writer = &fakeWriter{buffer: &buffer}
		fake.Closer = &fakeCloser{}

		connection := pop3.NewConnection(fake)

		response, err := connection.Cmd(tt.command)

		assert.Nil(t, err)
		assert.NotEmpty(t, t, response)
		assert.Equal(t, fmt.Sprintf("%s\r\n", tt.command), string(buffer))
	}
}
