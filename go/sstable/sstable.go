// Package sstable provides functionality to read and write SSTables.
// SSTables are sorted string tables, which are persistent, ordered, immutable
// maps from keys to values. This implementation is compatible with
// the SSTable format defined by https://github.com/mariusaeriksen/sstable.
package sstable

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

// errorCursor is a cursor that always returns an error.
type errorCursor struct {
	err error
}

func (c *errorCursor) Next() {}

func (c *errorCursor) Done() bool {
	return true
}

func (c *errorCursor) Error() error {
	return c.err
}

func (c *errorCursor) Entry() Entry {
	return Entry{}
}

// SSTable provides read-only access to an SSTable.
// It allows scanning entries from a given key or from the beginning.
// The underlying data is read from an io.Reader, io.ReaderAt, or io.ReadSeeker.
type SSTable struct {
	header header      // The SSTable header, contains metadata like index offset.
	index  index       // The SSTable index, used to find data blocks.
	reader interface{} // The underlying reader for SSTable data (e.g., a file).
	// It can be io.Reader, io.ReaderAt, or io.ReadSeeker.

	noCursor bool // noCursor is true if ScanFrom has been called on an io.Reader (non-seekable).
	// This ensures that scanning can only occur once for such readers, as they cannot be reset.
}

// NewSSTable creates a new SSTable reader from the given input source r.
// The input r can be an io.Reader, io.ReaderAt, or io.ReadSeeker.
//
// If r is an io.ReadSeeker or io.ReaderAt, the header and index are read
// immediately, allowing for random access scans.
// If r is only an io.Reader, the header is read, but the index cannot be
// fully loaded, meaning only a single, full sequential scan is supported.
//
// Returns an initialized SSTable or an error if the SSTable format is invalid
// or the reader operations fail.
func NewSSTable(r interface{}) (*SSTable, error) {
	table := SSTable{
		header: header{},
		index:  index{},
		reader: r, // Store the original reader.
	}

	// Handle different types of readers.
	switch r := r.(type) {
	case io.ReadSeeker:
		// For seekable readers, we can read the header and then the index.
		// Ensure we are at the beginning of the reader.
		newOffset, err := r.Seek(0, 0)
		if err != nil {
			return nil, fmt.Errorf("NewSSTable: failed to seek to beginning of reader: %w", err)
		}

		if newOffset != 0 {
			return nil, fmt.Errorf("NewSSTable: expected offset 0 after seeking to beginning, got %d", newOffset)
		}

		if err := table.header.read(r); err != nil {
			return nil, fmt.Errorf("NewSSTable: failed to read header: %w", err)
		}

		if int64(table.header.indexOffset) < 0 {
			return nil, fmt.Errorf("NewSSTable: invalid negative index offset %d", table.header.indexOffset)
		}

		newOffset, err = r.Seek(int64(table.header.indexOffset), 0)
		if err != nil {
			return nil, fmt.Errorf("NewSSTable: failed to seek to index offset %d: %w", table.header.indexOffset, err)
		}

		if uint64(newOffset) != table.header.indexOffset {
			return nil, fmt.Errorf("NewSSTable: expected offset %d after seeking to index, got %d", table.header.indexOffset, newOffset)
		}

		if _, err := table.index.ReadFrom(r); err != nil {
			// If reading the index fails, return the error.
			return nil, fmt.Errorf("NewSSTable: failed to read index from ReadSeeker: %w", err)
		}
	case io.ReaderAt:
		// For readers that support ReadAt, we can read the header from offset 0
		// and then the index from the offset specified in the header.
		headerBytes := make([]byte, headerSize)
		if n, err := r.ReadAt(headerBytes, 0); n != len(headerBytes) || err != nil {
			if err == nil && n != len(headerBytes) { // Check if err is nil before creating a new error
				// This case handles io.ReaderAt implementations that might not return an error on short reads.
				err = fmt.Errorf("NewSSTable: read %d bytes for header, expected %d", n, headerSize)
			}
			return nil, fmt.Errorf("NewSSTable: failed to read header bytes at offset 0: %w", err)
		}

		if err := table.header.UnmarshalBinary(headerBytes); err != nil {
			return nil, fmt.Errorf("NewSSTable: failed to unmarshal header: %w", err)
		}

		// Read the index from the specified offset.
		if err := table.index.ReadAt(r, table.header.indexOffset); err != nil {
			return nil, fmt.Errorf("NewSSTable: failed to read index at offset %d: %w", table.header.indexOffset, err)
		}
	case io.Reader:
		// For a simple io.Reader, we can only read the header.
		// The index cannot be loaded because we cannot seek.
		// This means only a single sequential scan from the beginning is possible.
		if err := table.header.read(r); err != nil {
			return nil, fmt.Errorf("NewSSTable: failed to read header for non-seekable reader: %w", err)
		}
		// The index remains empty. ScanFrom will handle this.
	default:
		// The provided reader type is not supported.
		return nil, fmt.Errorf("NewSSTable: unsupported reader type %T", r)
	}

	return &table, nil
}

// ScanFrom returns a new Cursor to iterate over entries in the SSTable.
//
// The scan starts from the first entry whose key is greater than or equal to
// the given `key`. If `key` is nil or empty, the scan starts from the
// beginning of the SSTable.
//
// The behavior depends on the type of the underlying reader:
// - io.ReaderAt or io.ReadSeeker: Supports multiple scans and seeking to `key`.
// - io.Reader: Supports only a single, sequential scan from the beginning.
//   Attempting a second scan or a scan with a non-nil `key` (if the index
//   isn't available) will result in a cursor that returns an error.
func (s *SSTable) ScanFrom(key []byte) Cursor {
	switch r := s.reader.(type) {
	case io.ReaderAt:
		// For io.ReaderAt, we can use the index to find the starting block.
		var startOffset uint64
		if key == nil {
			// If key is nil, start scan from the beginning of data blocks.
			startOffset = headerSize
		} else {
			// Find the index entry for the key.
			// entryIndexOf returns the index of the entry whose key is *just before*
			// or equal to the target key, or -1 if all keys are greater or key is smallest.
			// If the key is smaller than all keys in the index, i will be -1.
			// If the key is larger than all keys, i will be len(s.index)-1.
			i := s.index.entryIndexOf(key)
			if i == -1 {
				// Key is smaller than any indexed key, or index is empty.
				// Start from the very first data block if index was empty or key is before first indexed key.
				// Or, if index is not empty, this implies key is smaller than s.index[0].keyBytes
				// so we should start at s.index[0].blockOffset.
				// However, entryIndexOf returns -1 also if index is empty.
				// A safer bet if index is present is to use the first block.
				// If index is empty, headerSize is the only logical start.
				if len(s.index) > 0 {
					startOffset = s.index[0].blockOffset
				} else {
					startOffset = headerSize // No index, start from after header
				}
			} else {
				// Start from the block offset indicated by the index.
				startOffset = s.index[i].blockOffset
			}
		}

		c := CursorToOffset{
			reader:    r, // Use the type-asserted reader
			offset:    startOffset,
			endOffset: s.header.indexOffset, // Scan up to the start of the index.
		}

		// If a key is provided, advance the cursor until Entry().Key >= key.
		// This is necessary because the index points to blocks, not exact key locations.
		if key != nil {
			for !c.Done() && bytes.Compare(c.Entry().Key, key) < 0 {
				c.Next()
				if c.Error() != nil { // Check for errors during Next()
					return &errorCursor{err: c.Error()}
				}
			}
		}
		return &c
	case io.Reader:
		// For a simple io.Reader (which is not an io.ReadSeeker or io.ReaderAt),
		// we can only scan the data once because we cannot seek backwards.
		// The noCursor flag tracks if ScanFrom has already been called.
		if s.noCursor {
			return &errorCursor{err: fmt.Errorf("ScanFrom: cursor already obtained for non-seekable reader type %T", r)}
		}
		s.noCursor = true

		// With a simple io.Reader, we also don't have a loaded index if it wasn't an io.ReadSeeker.
		// So, we can only scan from the beginning.
		// If a key is provided, we try to advance, but this is inefficient without an index.
		if key != nil && len(s.index) == 0 { // len(s.index) == 0 implies it's a non-seekable io.Reader where index wasn't loaded
			// It's not ideal to scan from beginning for a key with plain io.Reader without index.
			// However, the original logic attempts this.
			// A better approach might be to return an error if key is not nil and it's a plain io.Reader.
			// For now, maintaining existing behavior.
		}

		c := CursorToOffset{
			reader:    r, // Use the type-asserted reader
			offset:    headerSize,             // Start from after the header.
			endOffset: s.header.indexOffset, // This might be 0 if header could not be fully processed for indexOffset.
			// For a pure io.Reader, indexOffset in header might not be known if reading index was skipped.
			// However, CursorToOffset should handle reaching EOF correctly.
			// If indexOffset is 0 from header read of a non-seekable reader, it means scan till EOF.
		}

		// If a key is provided, advance the cursor. This will read sequentially.
		if key != nil {
			for !c.Done() && bytes.Compare(c.Entry().Key, key) < 0 {
				c.Next()
				if c.Error() != nil { // Check for errors during Next()
					return &errorCursor{err: c.Error()}
				}
			}
		}
		return &c
	default:
		// Unsupported reader type.
		return &errorCursor{err: fmt.Errorf("ScanFrom: unsupported reader type %T", s.reader)}
	}
}
