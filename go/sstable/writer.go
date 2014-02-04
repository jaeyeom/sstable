package sstable

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// Writer is used to build a SSTable binary with Write function.
type Writer struct {
	indexBuffer indexBuffer
	lastKey     []byte
	writer      io.Writer
	closed      bool
}

// NewWriter creates a Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{
		indexBuffer: indexBuffer{
			maxBlockLength: 64 * 1024,
			offset:         uint64(0),
			index:          index{},
		},
		writer: w,
	}
}

// Write writes an entry. Multiple calls to the function appends
// entries to the SSTable. The call should be made in sorted order of
// the keys.
func (w *Writer) Write(e Entry) error {
	if w.indexBuffer.offset == 0 {
		h := header{2, 0, 0}
		offset, err := h.WriteTo(w.writer)
		if err != nil {
			return err
		}
		w.indexBuffer.offset = uint64(offset)
	}
	if w.lastKey != nil && bytes.Compare(w.lastKey, e.Key) > 0 {
		return fmt.Errorf("key is not sorted")
	}
	w.indexBuffer.Write(e.Key, uint32(len(e.Value)))
	_, err := e.WriteTo(w.writer)
	w.lastKey = e.Key
	return err
}

// Close closes the writer. It writes index at the end and overwrite
// header at front.
func (w *Writer) Close() error {
	if w.closed {
		return errors.New("Writer.Close: already closed")
	}
	_, err := w.indexBuffer.index.WriteTo(w.writer)
	h := header{
		version:     2,
		numBlocks:   uint32(len(w.indexBuffer.index)),
		indexOffset: w.indexBuffer.offset,
	}
	switch writer := w.writer.(type) {
	case io.WriterAt:
		data, err := h.MarshalBinary()
		if err != nil {
			return err
		}
		_, err = writer.WriteAt(data, 0)
		return err
	case io.WriteSeeker:
		writer.Seek(0, 0)
		_, err := h.WriteTo(writer)
		return err
	default:
		return errors.New("Writer.Close: writer cannot do random access")
	}
	if closer, ok := w.writer.(io.Closer); ok {
		closer.Close()
	}
	w.closed = true
	return err
}
