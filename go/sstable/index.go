package sstable

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"sort"
)

// indexEntry is an entry of index.
type indexEntry struct {
	// KeyLength   uint32
	blockOffset uint64
	blockLength uint32
	keyBytes    []byte
}

// readIndexEntry reads indexEntry from r and returns indexEntry,
// number of bytes read, and error.
func readIndexEntry(r io.Reader) (*indexEntry, int, error) {
	lenbuf := make([]byte, 4)
	n, err := io.ReadFull(r, lenbuf)
	if err != nil {
		return nil, n, err
	}
	length := binary.BigEndian.Uint32(lenbuf)
	buf := make([]byte, length+16)
	copy(buf[:4], lenbuf)
	nn, err := io.ReadFull(r, buf[4:])
	n += nn
	if err != nil {
		return nil, n, err
	}
	e := indexEntry{}
	return &e, n, e.UnmarshalBinary(buf)
}

// readIndexEntryAt reads indexEntry at offset from r and returns
// indexEntry and error.
func readIndexEntryAt(r io.ReaderAt, offset uint64) (*indexEntry, error) {
	if int64(offset) < 0 {
		panic("unimplemented")
	}
	lenbuf := make([]byte, 4)
	if n, err := r.ReadAt(lenbuf, int64(offset)); n != len(lenbuf) {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lenbuf)
	buf := make([]byte, length+16)
	if n, err := r.ReadAt(buf, int64(offset)); n != len(buf) {
		return nil, err
	}
	e := indexEntry{}
	return &e, e.UnmarshalBinary(buf)
}

// size returns number of bytes of the indexEntry.
func (e *indexEntry) size() int {
	return 16 + len(e.keyBytes)
}

// WriteTo implements the io.WriterTo interface.
func (e *indexEntry) WriteTo(w io.Writer) (n int64, err error) {
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
func (e *indexEntry) MarshalBinary() (data []byte, err error) {
	buf := bytes.NewBuffer([]byte{})
	if err = binary.Write(buf, binary.BigEndian, uint32(len(e.keyBytes))); err != nil {
		return
	}
	if err = binary.Write(buf, binary.BigEndian, e.blockOffset); err != nil {
		return
	}
	if err = binary.Write(buf, binary.BigEndian, e.blockLength); err != nil {
		return
	}
	if _, err = buf.Write(e.keyBytes); err != nil {
		return
	}
	data = buf.Bytes()
	return
}

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (e *indexEntry) UnmarshalBinary(data []byte) error {
	if len(data) < 16 {
		return errors.New("indexEntry.UnmarshalBinary: invalid length")
	}
	length := binary.BigEndian.Uint32(data[:4])
	blockOffset := binary.BigEndian.Uint64(data[4:12])
	blockLength := binary.BigEndian.Uint32(data[12:16])
	if length != uint32(len(data[16:])) {
		return errors.New("indexEntry.UnmarshalBinary: invalid length")
	}
	e.blockOffset = blockOffset
	e.blockLength = blockLength
	e.keyBytes = make([]byte, length)
	copy(e.keyBytes, data[16:])
	return nil
}

// index stores index entries and implements a find function.
type index []indexEntry

// entryIndexOf returns the index of index entry that might contain
// the key. It returns -1 if there is no index entry.
func (i index) entryIndexOf(key []byte) int {
	return sort.Search(len(i), func(idx int) bool {
		return bytes.Compare(i[idx].keyBytes, key) > 0
	}) - 1
}

// ReadFrom implements the io.ReaderFrom interface.
func (i *index) ReadFrom(r io.Reader) (n int64, err error) {
	var e *indexEntry
	var nn int
	for err == nil {
		e, nn, err = readIndexEntry(r)
		n += int64(nn)
		if err == nil || err == io.EOF && nn > 0 {
			*i = append(*i, *e)
			continue
		}
		if err == io.EOF {
			err = nil
		}
		return
	}
	panic("unreachable")
}

// ReadAt reads index from r at offset.
func (i *index) ReadAt(r io.ReaderAt, offset uint64) (err error) {
	var e *indexEntry
	for err == nil {
		e, err = readIndexEntryAt(r, offset)
		if e != nil {
			*i = append(*i, *e)
			offset += uint64(e.size())
			continue
		}
		if err == io.EOF {
			err = nil
		}
		return
	}
	panic("unreachable")
}

// WriteTo implements the io.WriterTo interface.
func (i index) WriteTo(w io.Writer) (n int64, err error) {
	var nn int64
	for _, entry := range i {
		nn, err = entry.WriteTo(w)
		n += nn
		if err != nil {
			return
		}
	}
	return
}

// indexBuffer implements functions to build a new index.
type indexBuffer struct {
	maxBlockLength uint32
	offset         uint64
	index          index
}

// Write writes an entry in the buffer to build the index.
func (w *indexBuffer) Write(key []byte, valueSize uint32) {
	size := len(w.index)
	if size == 0 || int64(w.index[size-1].blockLength)+int64(valueSize) > int64(w.maxBlockLength) {
		w.index = append(w.index, indexEntry{})
		size += 1
		w.index[size-1].blockOffset = w.offset
		w.index[size-1].keyBytes = key
	}
	w.offset += uint64(8) + uint64(len(key)) + uint64(valueSize)
	w.index[size-1].blockLength += uint32(8) + uint32(len(key)) + valueSize
}
