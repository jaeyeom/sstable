package sstable

import (
	"encoding/hex"
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
	fmt.Print(hex.Dump(b))
	// Output:
	// 00000000  00 00 00 01 00 00 00 02  00 00 00 00 00 00 00 03  |................|
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
