package sstable

import (
	"fmt"
)

//nolint:govet
func Example_headerMarshalBinary() {
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

//nolint:govet
func Example_headerUnmarshalBinary() {
	h := header{}

	err := h.UnmarshalBinary([]byte{
		0, 0, 0, 1,
		0, 0, 0, 2,
		0, 0, 0, 0, 0, 0, 0, 3,
	})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(h)
	// Output:
	// {1 2 3}
}
