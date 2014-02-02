package sstable

import (
	"fmt"
)

func ExampleHeaderMarshalBinary() {
	h := header{
		version:     1,
		numBlocks:   2,
		indexOffset: 3,
	}
	b, _ := h.MarshalBinary()
	fmt.Println(b)
	// Output:
	// [0 0 0 1 0 0 0 2 0 0 0 0 0 0 0 3]
}

func ExampleHeaderUnmarshalBinary() {
	h := header{}
	h.UnmarshalBinary([]byte{0, 0, 0, 1, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 3})
	fmt.Println(h)
	// Output:
	// {1 2 3}
}
