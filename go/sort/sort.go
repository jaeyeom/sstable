// Package sort implements external merge sort.
package sort

import (
	"container/heap"

	"github.com/jaeyeom/sstable/go/sstable"
)

// SortEntries sorts all entries from r in memory and write to w.
func SortEntries(c sstable.Cursor, w *sstable.Writer) error {
	es := Entries{}
	for !c.Done() {
		e := c.Entry()
		c.Next()
		es = append(es, HeapEntry{*e, nil})
	}
	heap.Init(&es)
	for es.Len() > 0 {
		e := heap.Pop(&es)
		w.Write(e.(HeapEntry).Entry)
	}
	return w.Close()
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
