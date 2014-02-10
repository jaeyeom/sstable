package sstable

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

// Entry struct is a key value pair.
type Entry struct {
	// KeyLength   uint32
	// ValueLength uint32
	Key   []byte
	Value []byte
}

// readEntry reads an entry from r.
func readEntry(r io.Reader) (*Entry, error) {
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
	e := Entry{}
	return &e, e.UnmarshalBinary(buf)
}

// readEntryAt reads an entry from the offset of r.
func readEntryAt(r io.ReaderAt, offset uint64) (*Entry, error) {
	lenbuf := make([]byte, 8)
	if int64(offset) < 0 {
		panic("unimplemented")
	}
	if n, err := r.ReadAt(lenbuf, int64(offset)); n != len(lenbuf) {
		return nil, err
	}
	keyLength := binary.BigEndian.Uint32(lenbuf[:4])
	valueLength := binary.BigEndian.Uint32(lenbuf[4:8])
	buf := make([]byte, 8+keyLength+valueLength)
	if n, err := r.ReadAt(buf, int64(offset)); n != len(buf) {
		return nil, err
	}
	e := Entry{}
	return &e, e.UnmarshalBinary(buf)
}

// size returns number of bytes in this entry.
func (e *Entry) size() int {
	return 8 + len(e.Key) + len(e.Value)
}

// WriteTo implements the io.WriterTo interface.
func (e *Entry) WriteTo(w io.Writer) (n int64, err error) {
	var data []byte
	data, err = e.MarshalBinary()
	if err != nil {
		return
	}
	nn, err := w.Write(data)
	n = int64(nn)
	return
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (e *Entry) MarshalBinary() (data []byte, err error) {
	buf := bytes.NewBuffer([]byte{})
	if err = binary.Write(buf, binary.BigEndian, uint32(len(e.Key))); err != nil {
		return
	}
	if err = binary.Write(buf, binary.BigEndian, uint32(len(e.Value))); err != nil {
		return
	}
	if _, err = buf.Write(e.Key); err != nil {
		return
	}
	if _, err = buf.Write(e.Value); err != nil {
		return
	}
	data = buf.Bytes()
	return
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (e *Entry) UnmarshalBinary(data []byte) error {
	if len(data) < 8 {
		return errors.New("Entry.UnmarshalBinary: invalid length")
	}
	keyLength := binary.BigEndian.Uint32(data[:4])
	valueLength := binary.BigEndian.Uint32(data[4:8])
	if uint32(len(data)) != uint32(8)+keyLength+valueLength {
		return errors.New("Entry.UnmarshalBinary: invalid length")
	}
	e.Key = make([]byte, keyLength)
	copy(e.Key, data[8:8+keyLength])
	e.Value = make([]byte, valueLength)
	copy(e.Value, data[8+keyLength:])
	return nil
}
