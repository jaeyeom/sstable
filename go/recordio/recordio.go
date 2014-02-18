// Package recordio implements Closer interface to recordio.Writer.
package recordio

import (
	"io"

	"github.com/eclesh/recordio"
)

// WriteCloser implements Closer interface to recordio.Writer.
type WriteCloser struct {
	recordio.Writer
	f io.Writer
}

// NewWriteCloser returns a new WriteCloser.
func NewWriteCloser(w io.Writer) *WriteCloser {
	return &WriteCloser{*recordio.NewWriter(w), w}
}

// Close closes the file.
func (w *WriteCloser) Close() error {
	if c, ok := interface{}(&w.Writer).(io.Closer); ok {
		return c.Close()
	}
	if c, ok := w.f.(io.Closer); ok {
		return c.Close()
	}
	return nil
}
