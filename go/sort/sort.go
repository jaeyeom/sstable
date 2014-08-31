// Package sort implements external merge sort.
package sort

import (
	"container/heap"

	"github.com/jaeyeom/sstable/go/sstable"
)

// SortEntries sorts entries from c in memory and write to w in
// slightly over the maxSize bytes. Returns number of entries written
// and error.
func SortEntries(c sstable.Cursor, maxSize uint64, w *sstable.Writer) (n int, err error) {
	es := Entries{}
	size := uint64(0)
	for !c.Done() && size < maxSize {
		e := c.Entry()
		c.Next()
		size += e.Size()
		es = append(es, HeapEntry{*e, nil})
	}
	heap.Init(&es)
	for es.Len() > 0 {
		e := heap.Pop(&es)
		err = w.Write(e.(HeapEntry).Entry)
		if err != nil {
			return
		}
		n += 1
	}
	err = w.Close()
	return
}

// Merge merges from multiple cursors and write SSTable to w.
func Merge(cursors []sstable.Cursor, w *sstable.Writer) error {
	es := Entries{}
	for i, c := range cursors {
		if c.Done() {
			continue
		}
		e := c.Entry()
		c.Next()
		es = append(es, HeapEntry{*e, i})
	}
	heap.Init(&es)
	for es.Len() > 0 {
		e := heap.Pop(&es).(HeapEntry)
		i := e.data.(int)
		if !cursors[i].Done() {
			heap.Push(&es, HeapEntry{*cursors[i].Entry(), i})
			cursors[i].Next()
		}
		if err := w.Write(e.Entry); err != nil {
			return err
		}
	}
	return nil
}
