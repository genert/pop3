package pop3_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/genert/pop3"
)

func TestIsOK(t *testing.T) {
	testCases := []struct {
		message        string
		expectedResult bool
	}{
		{"+OK", true},
		{"+OK 2 messages", true},
		{"-ERR", false},
		{"-ERR no such message", false},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf(`IsOK should return %s for "%s" message`, strconv.FormatBool(tt.expectedResult), tt.message), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expectedResult, pop3.IsOK(tt.message))
		})
	}
}

func TestIsErr(t *testing.T) {
	testCases := []struct {
		message        string
		expectedResult bool
	}{
		{"+OK", false},
		{"+OK 2 messages", false},
		{"-ERR", true},
		{"-ERR no such message", true},
	}

	for _, tt := range testCases {
		t.Run(fmt.Sprintf(`IsErr should return %s for "%s" message`, strconv.FormatBool(tt.expectedResult), tt.message), func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tt.expectedResult, pop3.IsErr(tt.message))
		})
	}
}
