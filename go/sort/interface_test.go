package sort

import (
	"container/heap"
	"fmt"

	"github.com/jaeyeom/sstable/go/sstable"
)

func ExampleHeap() {
	h := &Entries{}
	heap.Push(h, sstable.Entry{[]byte("key3"), []byte("value")})
	heap.Push(h, sstable.Entry{[]byte("key1"), []byte("value")})
	heap.Push(h, sstable.Entry{[]byte("key4"), []byte("value")})
	heap.Push(h, sstable.Entry{[]byte("key2"), []byte("value2")})
	heap.Push(h, sstable.Entry{[]byte("key2"), []byte("value")})
	fmt.Println(heap.Pop(h))
	fmt.Println(heap.Pop(h))
	fmt.Println(heap.Pop(h))
	fmt.Println(heap.Pop(h))
	fmt.Println(heap.Pop(h))
	// Output:
	// {[107 101 121 49] [118 97 108 117 101]}
	// {[107 101 121 50] [118 97 108 117 101]}
	// {[107 101 121 50] [118 97 108 117 101 50]}
	// {[107 101 121 51] [118 97 108 117 101]}
	// {[107 101 121 52] [118 97 108 117 101]}
}
