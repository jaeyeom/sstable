// Package shard implements sharded writer.
package shard

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"hash"
	"os"
)

// PrefixSum64Hash implements hash.Hash64 interface.
type PrefixSum64Hash struct {
	hash.Hash
}

// Sun64 implements hash.Hash64 interface. It simply reads the first 8
// bytes by big endian.
func (h *PrefixSum64Hash) Sum64() uint64 {
	var sum uint64
	b := h.Sum([]byte{})
	buf := bytes.NewReader(b)
	err := binary.Read(buf, binary.BigEndian, &sum)
	if err != nil {
		// Should never happen.
		panic(err)
	}
	return sum
}

// Writer writes data to one of the shard based on the hash function.
type Writer struct {
	w []io.Writer
	h hash.Hash64
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
func NewWriter(n int, h hash.Hash64, wf WriterFactory) (*Writer, error) {
	w := Writer{
		w: make([]io.Writer, n),
		h: h,
	}
	for i := 0; i < n; i++ {
		w.w[i] = wf(i, n)
	}
	return &w, nil
}

// Write writes data to sharded writer.
func (w *Writer) Write(data []byte) (int, error) {
	w.h.Reset()
	w.h.Write(data)
	i := w.h.Sum64() % uint64(len(w.w))
	return w.w[i].Write(data)
}

// Close closes shareded writer and never return error.
func (w *Writer) Close() error {
	for i := 0; i < len(w.w); i++ {
		if c, ok := w.w[i].(io.Closer); ok {
			c.Close()
		}
	}
	return nil
}
