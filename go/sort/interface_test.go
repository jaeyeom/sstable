package sort

import (
	"container/heap"
	"fmt"

	"github.com/jaeyeom/sstable/go/sstable"
)

func Example_heap() {
	h := &Entries{}
	heap.Push(h, HeapEntry{sstable.Entry{
		Key:   []byte("key3"),
		Value: []byte("value"),
	}, 1})
	heap.Push(h, HeapEntry{sstable.Entry{
		Key:   []byte("key1"),
		Value: []byte("value"),
	}, 2})
	heap.Push(h, HeapEntry{sstable.Entry{
		Key:   []byte("key4"),
		Value: []byte("value"),
	}, 3})
	heap.Push(h, HeapEntry{sstable.Entry{
		Key:   []byte("key2"),
		Value: []byte("value2"),
	}, 4})
	heap.Push(h, HeapEntry{sstable.Entry{
		Key:   []byte("key2"),
		Value: []byte("value"),
	}, 5})
	fmt.Println(heap.Pop(h))
	fmt.Println(heap.Pop(h))
	fmt.Println(heap.Pop(h))
	fmt.Println(heap.Pop(h))
	fmt.Println(heap.Pop(h))
	// Output:
	// {{[107 101 121 49] [118 97 108 117 101]} 2}
	// {{[107 101 121 50] [118 97 108 117 101]} 5}
	// {{[107 101 121 50] [118 97 108 117 101 50]} 4}
	// {{[107 101 121 51] [118 97 108 117 101]} 1}
	// {{[107 101 121 52] [118 97 108 117 101]} 3}
}
