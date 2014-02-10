package sstable

import (
	"bytes"
	"errors"
	"io"
)

// SSTable implements read only random access of the SSTable.
type SSTable struct {
	header header
	index  index
	reader interface{}
}

// NewSSTable creates a SSTable struct
func NewSSTable(r interface{}) (*SSTable, error) {
	table := SSTable{
		header: header{},
		index:  index{},
		reader: r,
	}
	switch r := r.(type) {
	case io.ReadSeeker:
		newOffset, err := r.Seek(0, 0)
		if err != nil {
			return nil, err
		}
		if newOffset != 0 {
			return nil, errors.New("NewSSTable: new offset is not zero")
		}
		if err := table.header.read(r); err != nil {
			return nil, err
		}
		if int64(table.header.indexOffset) < 0 {
			panic("unimplemented")
		}
		newOffset, err = r.Seek(int64(table.header.indexOffset), 0)
		if err != nil {
			return nil, err
		}
		if uint64(newOffset) != table.header.indexOffset {
			return nil, errors.New("NewSSTable: new offset is not same as the index offset")
		}
		if _, err := table.index.ReadFrom(r); err != nil {
			return nil, err
		}
	case io.ReaderAt:
		headerBytes := make([]byte, headerSize)
		if n, err := r.ReadAt(headerBytes, 0); n != len(headerBytes) {
			return nil, err
		}
		if err := table.header.UnmarshalBinary(headerBytes); err != nil {
			return nil, err
		}
		if err := table.index.ReadAt(r, table.header.indexOffset); err != nil {
			return nil, err
		}
	default:
		panic("unimplemented")
	}
	return &table, nil
}

// block returns ith block of the SSTable.
func (s *SSTable) block(i int) ([]byte, error) {
	switch r := s.reader.(type) {
	case io.ReaderAt:
		b := make([]byte, s.index[i].blockLength)
		if int64(s.index[i].blockOffset) < 0 {
			panic("unimplemented")
		}
		if int(s.index[i].blockLength) < 0 {
			panic("unimplemented")
		}
		if n, err := r.ReadAt(b, int64(s.index[i].blockOffset)); n != int(s.index[i].blockLength) {
			return nil, err
		}
		return b, nil
	default:
		panic("unimplemented")
	}
}

// ScanFrom scans from the key to the end of the SSTable.
func (s *SSTable) ScanFrom(key []byte) Cursor {
	i := s.index.entryIndexOf(key)
	if i == -1 {
		i = 0
	}
	c := CursorToOffset{
		table:     s,
		offset:    s.index[i].blockOffset,
		endOffset: s.header.indexOffset,
	}
	for !c.Done() && bytes.Compare(c.Entry().Key, key) < 0 {
		c.Next()
	}
	return &c
}
