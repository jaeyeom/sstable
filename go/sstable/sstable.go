package sstable

import (
	"bytes"
	"errors"
	"io"
)

// SSTable implements read only random access of the SSTable.
type SSTable struct {
	header   header
	index    index
	reader   interface{}
	noCursor bool
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
	case io.Reader:
		// Index can't be read if the reader isn't random access.
		if err := table.header.read(r); err != nil {
			return nil, err
		}
	default:
		panic("unimplemented")
	}

	return &table, nil
}

// ScanFrom scans from the key to the end of the SSTable. If key is
// nil, scan from the beginning.
func (s *SSTable) ScanFrom(key []byte) Cursor {
	switch s.reader.(type) {
	case io.ReaderAt:
		if key == nil {
			return &CursorToOffset{
				reader:    s.reader,
				offset:    headerSize,
				endOffset: s.header.indexOffset,
			}
		}

		i := s.index.entryIndexOf(key)
		if i == -1 {
			i = 0
		}

		c := CursorToOffset{
			reader:    s.reader,
			offset:    s.index[i].blockOffset,
			endOffset: s.header.indexOffset,
		}

		if key != nil {
			for !c.Done() && bytes.Compare(c.Entry().Key, key) < 0 {
				c.Next()
			}
		}

		return &c
	case io.Reader:
		if s.noCursor {
			panic("unimplemented")
		}

		s.noCursor = true

		c := CursorToOffset{
			reader:    s.reader,
			offset:    headerSize,
			endOffset: s.header.indexOffset,
		}

		if key != nil {
			for !c.Done() && bytes.Compare(c.Entry().Key, key) < 0 {
				c.Next()
			}
		}

		return &c
	default:
		panic("unimplemented")
	}
}
