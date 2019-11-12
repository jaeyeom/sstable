package sstable

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"golang.org/x/xerrors"
)

// headerSize is the number of bytes of the header.
const headerSize = 16

// header implements binary IO and marshal functions.
type header struct {
	version     uint32
	numBlocks   uint32
	indexOffset uint64
}

// read reads and parses a header from r.
func (h *header) read(r io.Reader) error {
	var buf [headerSize]byte

	if _, err := io.ReadFull(r, buf[:]); err != nil {
		return xerrors.Errorf("failed to read the header: %w", err)
	}

	if err := h.UnmarshalBinary(buf[:]); err != nil {
		return xerrors.Errorf("failed to unmarshal header: %w", err)
	}

	return nil
}

// WriteTo implements the io.WriterTo interface.
func (h *header) WriteTo(w io.Writer) (n int64, err error) {
	data, err := h.MarshalBinary()
	if err != nil {
		return
	}

	nn, err := w.Write(data)
	return int64(nn), err //nolint:wsl
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (h *header) MarshalBinary() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0, headerSize))
	if err := binary.Write(buf, binary.BigEndian, h); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (h *header) UnmarshalBinary(data []byte) error {
	if len(data) != headerSize {
		return errors.New("header.UnmarshalBinary: invalid length")
	}

	h.version = binary.BigEndian.Uint32(data[:4])
	h.numBlocks = binary.BigEndian.Uint32(data[4:8])
	h.indexOffset = binary.BigEndian.Uint64(data[8:])

	return nil
}
