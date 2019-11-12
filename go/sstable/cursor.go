package sstable

import (
	"io"
)

// Cursor is an interface to iterate.
type Cursor interface {
	Entry() *Entry
	Done() bool
	Next()
}

// CursorToOffset is a Cursor that read until the endOffset.
type CursorToOffset struct {
	reader    interface{}
	offset    uint64
	endOffset uint64
	entry     *Entry
}

// Entry returns the current entry.
func (c *CursorToOffset) Entry() *Entry {
	if c.entry != nil {
		return c.entry
	}

	switch r := c.reader.(type) {
	case io.ReaderAt:
		e, err := ReadEntryAt(r, c.offset)
		if err != nil {
			return nil
		}

		c.offset += e.Size()
		c.entry = e

		return c.entry
	case io.Reader:
		e, err := ReadEntry(r)
		if err != nil {
			return nil
		}

		c.offset += e.Size()
		c.entry = e

		return c.entry
	default:
		panic("unimplemented")
	}
}

// Done returns true when there is no more entry to read.
func (c *CursorToOffset) Done() bool {
	if c.entry == nil && c.offset >= c.endOffset {
		c.reader = nil
		return true
	}

	return false
}

// Next moves the cursor to the next entry.
func (c *CursorToOffset) Next() {
	if c.entry == nil {
		c.Entry()
	}

	c.entry = nil
}
