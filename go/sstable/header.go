package sstable

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// header implements binary IO and marshal functions.
type header struct {
	version     uint32
	numBlocks   uint32
	indexOffset uint64
}

// readHeader reads and parses a header from r.
func readHeader(r io.Reader) (*header, error) {
	var buf [16]byte
	if _, err := io.ReadFull(r, buf[:16]); err != nil {
		return nil, err
	}
	h := header{}
	h.UnmarshalBinary(buf[:])
	return &h, nil
}

// WriteTo implements the io.WriterTo interface.
func (h *header) WriteTo(w io.Writer) (n int64, err error) {
	var data []byte
	data, err = h.MarshalBinary()
	if err != nil {
		return
	}
	nn, err := w.Write(data)
	n = int64(nn)
	return
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (h *header) MarshalBinary() (data []byte, err error) {
	buf := bytes.NewBuffer(make([]byte, 0, 16))
	if err := binary.Write(buf, binary.BigEndian, h); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (h *header) UnmarshalBinary(data []byte) error {
	if len(data) != 16 {
		return errors.New("header.UnmarshalBinary: invalid length")
	}
	h.version = binary.BigEndian.Uint32(data[:4])
	h.numBlocks = binary.BigEndian.Uint32(data[4:8])
	h.indexOffset = binary.BigEndian.Uint64(data[8:])
	return nil
}
