// Package shard implements sharded writer.
package shard

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash"
	"io"
	"os"
)

// PrefixSum64Hash implements hash.Hash64 interface.
type PrefixSum64Hash struct {
	hash.Hash
}

// Sun64 implements hash.Hash64 interface. It simply reads the first 8
// bytes by big endian.
func (h *PrefixSum64Hash) Sum64() uint64 {
	b := h.Sum(nil)
	buf := bytes.NewReader(b)

	var sum uint64
	if err := binary.Read(buf, binary.BigEndian, &sum); err != nil {
		panic(err) // Should never happen.
	}

	return sum
}

type hash64 interface {
	io.Writer
	// Reset resets the Hash to its initial state.
	Reset()
	Sum64() uint64
}

// Writer writes data to one of the shard based on the hash function.
type Writer struct {
	w []io.Writer
	h hash64
}

// WriterFactory returns a writer for each i of n.
type WriterFactory func(i, n int) io.Writer

// NewOSFileWriterFactory returns a WriterFactory that returns os.File
// with the filenane prefixed by prefix and the index string such as
// 00000-of-00100, 00001-of-00100, and so on.
func NewOSFileWriterFactory(prefix string) WriterFactory {
	return func(i, n int) io.Writer {
		f, _ := os.Create(fmt.Sprintf("%s%05d-of-%05d", prefix, i, n))
		return f
	}
}

// NewWriter returns a new sharded writer.
func NewWriter(n int, h hash64, wf WriterFactory) *Writer {
	w := Writer{
		w: make([]io.Writer, n),
		h: h,
	}
	for i := 0; i < n; i++ {
		w.w[i] = wf(i, n)
	}

	return &w
}

// Write writes data to sharded writer.
func (w *Writer) Write(data []byte) (int, error) {
	w.h.Reset()

	if _, err := w.h.Write(data); err != nil {
		return 0, fmt.Errorf("hash write failed: %w", err)
	}

	i := w.h.Sum64() % uint64(len(w.w))

	return w.w[i].Write(data)
}

// Close closes sharded writer and never return error.
func (w *Writer) Close() error {
	for i := 0; i < len(w.w); i++ {
		if c, ok := w.w[i].(io.Closer); ok {
			c.Close()
		}
	}

	return nil
}
