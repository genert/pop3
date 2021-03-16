package pop3_test

import (
	"io"
)

type readWriteFaker struct {
	io.Reader
	io.Writer
	io.Closer
}

type fakeCloser struct{}

func (f fakeCloser) Close() error { return nil }

type fakeWriter struct {
	buffer *[]byte
}

func (f fakeWriter) Write(p []byte) (int, error) {
	*f.buffer = append(*f.buffer, p...)
	return len(p), nil
}

func (f fakeWriter) Flush() error { return nil }
