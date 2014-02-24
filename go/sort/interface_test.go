package sort

import (
	"container/heap"
	"fmt"

	"github.com/jaeyeom/sstable/go/sstable"
)

func ExampleHeap() {
	h := &Entries{}
	heap.Push(h, HeapEntry{sstable.Entry{[]byte("key3"), []byte("value")}, 1})
	heap.Push(h, HeapEntry{sstable.Entry{[]byte("key1"), []byte("value")}, 2})
	heap.Push(h, HeapEntry{sstable.Entry{[]byte("key4"), []byte("value")}, 3})
	heap.Push(h, HeapEntry{sstable.Entry{[]byte("key2"), []byte("value2")}, 4})
	heap.Push(h, HeapEntry{sstable.Entry{[]byte("key2"), []byte("value")}, 5})
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
