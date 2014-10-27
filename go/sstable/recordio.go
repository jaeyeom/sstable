package sstable

import (
	"io"
)

// NewRecordIOReader returns a cursor that reads RecordIO from r. It
// requires size.
func NewRecordIOReader(r io.Reader, size uint64) Cursor {
	return &CursorToOffset{
		reader:    r,
		offset:    0,
		endOffset: size,
	}
}
