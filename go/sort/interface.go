//nolint:revive // package name intentionally matches functionality
package sort

import (
	"bytes"

	"github.com/jaeyeom/sstable/go/sstable"
)

// HeapEntry is a struct with an entry and an additional data.
type HeapEntry struct {
	sstable.Entry
	data interface{}
}

// Entries implements the heap.Interface interface.
type Entries []HeapEntry

// Len implements the sort.Interface interface.
func (es Entries) Len() int {
	return len(es)
}

// Less implements the sort.Interface interface.
func (es Entries) Less(i, j int) bool {
	c := bytes.Compare(es[i].Key, es[j].Key)
	if c == 0 {
		return bytes.Compare(es[i].Value, es[j].Value) == -1
	}
	return c == -1 //nolint:wsl
}

// Swap implements the sort.Interface interface.
func (es Entries) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

// Push implements the heap.Interface interface.
func (es *Entries) Push(x interface{}) {
	if x, ok := x.(HeapEntry); ok {
		*es = append(*es, x)
	} else {
		panic("wrong type")
	}
}

// Pop implements the heap.Interface interface.
func (es *Entries) Pop() interface{} {
	last := (*es)[es.Len()-1]
	*es = (*es)[:es.Len()-1]

	return last
}
