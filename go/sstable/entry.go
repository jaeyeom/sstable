package sstable

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"math"
)

// Entry struct is a key value pair.
type Entry struct {
	// KeyLength   uint32
	// ValueLength uint32
	Key   []byte
	Value []byte
}

// ReadEntry reads an entry from r.
func ReadEntry(r io.Reader) (*Entry, error) {
	lenbuf := make([]byte, 8)
	if _, err := io.ReadFull(r, lenbuf); err != nil {
		return nil, err
	}

	keyLength := binary.BigEndian.Uint32(lenbuf[:4])
	valueLength := binary.BigEndian.Uint32(lenbuf[4:8])
	buf := make([]byte, 8+keyLength+valueLength)
	copy(buf, lenbuf)

	if _, err := io.ReadFull(r, buf[8:]); err != nil {
		return nil, err
	}

	var e Entry
	return &e, e.UnmarshalBinary(buf) //nolint:wsl
}

// ReadEntryAt reads an entry from the offset of r.
func ReadEntryAt(r io.ReaderAt, offset uint64) (*Entry, error) {
	var lenbuf [8]byte

	if offset > math.MaxInt64 {
		panic("unimplemented")
	}

	if n, err := r.ReadAt(lenbuf[:], int64(offset)); n != len(lenbuf) { //nolint:gosec // overflow checked above
		return nil, err
	}

	keyLength := binary.BigEndian.Uint32(lenbuf[:4])
	valueLength := binary.BigEndian.Uint32(lenbuf[4:8])
	buf := make([]byte, 8+keyLength+valueLength)

	if n, err := r.ReadAt(buf, int64(offset)); n != len(buf) { //nolint:gosec // overflow checked above
		return nil, err
	}

	var e Entry
	return &e, e.UnmarshalBinary(buf) //nolint:wsl
}

// Size returns number of bytes in this entry.
func (e *Entry) Size() uint64 {
	return uint64(8) + uint64(len(e.Key)) + uint64(len(e.Value))
}

// WriteTo implements the io.WriterTo interface.
func (e *Entry) WriteTo(w io.Writer) (n int64, err error) {
	data, err := e.MarshalBinary()
	if err != nil {
		return 0, err
	}

	nn, err := w.Write(data)
	return int64(nn), err //nolint:wsl
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (e *Entry) MarshalBinary() ([]byte, error) {
	if len(e.Key) > math.MaxUint32 || len(e.Value) > math.MaxUint32 {
		return nil, errors.New("Entry.MarshalBinary: key or value too large")
	}
	buf := bytes.NewBuffer([]byte{})
	if err := binary.Write(buf, binary.BigEndian, uint32(len(e.Key))); err != nil { //nolint:gosec // overflow checked above
		return nil, err
	}

	if err := binary.Write(buf, binary.BigEndian, uint32(len(e.Value))); err != nil { //nolint:gosec // overflow checked above
		return nil, err
	}

	if _, err := buf.Write(e.Key); err != nil {
		return nil, err
	}

	if _, err := buf.Write(e.Value); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (e *Entry) UnmarshalBinary(data []byte) error {
	if len(data) < 8 {
		return errors.New("Entry.UnmarshalBinary: invalid length")
	}

	keyLength, valueLength := binary.BigEndian.Uint32(data[:4]), binary.BigEndian.Uint32(data[4:8])
	expectedLen := uint64(8) + uint64(keyLength) + uint64(valueLength)
	if uint64(len(data)) != expectedLen {
		return errors.New("Entry.UnmarshalBinary: invalid length")
	}

	e.Key = make([]byte, keyLength)
	copy(e.Key, data[8:8+keyLength])
	e.Value = make([]byte, valueLength)
	copy(e.Value, data[8+keyLength:])

	return nil
}
