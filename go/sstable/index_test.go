package sstable

import (
	"bytes"
	"encoding/hex"
	"fmt"
)

//nolint:govet
func Example_indexBufferWrite() {
	w := &indexBuffer{
		maxBlockLength: 64 * 1024,
	}
	w.Write([]byte{1, 2, 3}, 30000)
	w.Write([]byte{1, 2, 3, 4}, 30000)
	w.Write([]byte{2, 3, 4}, 30000)
	fmt.Println(w.index)
	// Output:
	// [{0 60023 [1 2 3]} {60023 30011 [2 3 4]}]
}

//nolint:govet
func Example_indexEntryIndexOf() {
	i := &index{
		{0, 60023, []byte{1, 2, 3}},
		{60023, 30011, []byte{2, 3, 4}},
	}
	fmt.Println(i.entryIndexOf([]byte{1, 2}))
	fmt.Println(i.entryIndexOf([]byte{1, 2, 3}))
	fmt.Println(i.entryIndexOf([]byte{1, 2, 3, 4}))
	fmt.Println(i.entryIndexOf([]byte{2, 3, 4}))
	fmt.Println(i.entryIndexOf([]byte{2, 3, 5}))
	// Output:
	// -1
	// 0
	// 0
	// 1
	// 1
}

//nolint:govet
func Example_indexReadFrom() {
	var i index

	buf := bytes.NewBuffer([]byte{
		0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 234, 119, 1, 2, 3,
		0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 234, 119, 0, 0, 117, 59, 2, 3, 4,
	})
	if _, err := i.ReadFrom(buf); err != nil {
		fmt.Println(err)
	}

	fmt.Println(i)
	// Output:
	// [{0 60023 [1 2 3]} {60023 30011 [2 3 4]}]
}

//nolint:govet
func Example_indexReadAt() {
	var i index

	f := bytes.NewReader([]byte{
		0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 234, 119, 1, 2, 3,
		0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 234, 119, 0, 0, 117, 59, 2, 3, 4,
	})
	if err := i.ReadAt(f, 0); err != nil {
		fmt.Println(err)
	}

	fmt.Println(i)
	// Output:
	// [{0 60023 [1 2 3]} {60023 30011 [2 3 4]}]
}

//nolint:govet
func Example_indexWriteTo() {
	i := &index{
		{0, 60023, []byte{1, 2, 3}},
		{60023, 30011, []byte{2, 3, 4}},
	}

	buf := bytes.NewBuffer([]byte{})
	if _, err := i.WriteTo(buf); err != nil {
		fmt.Println(err)
	}

	fmt.Println(hex.Dump(buf.Bytes()))
	// Output:
	// 00000000  00 00 00 03 00 00 00 00  00 00 00 00 00 00 ea 77  |...............w|
	// 00000010  01 02 03 00 00 00 03 00  00 00 00 00 00 ea 77 00  |..............w.|
	// 00000020  00 75 3b 02 03 04                                 |.u;...|
}

//nolint:govet
func Example_indexEntryMarshalBinary() {
	b := indexEntry{
		blockOffset: 1,
		blockLength: 10,
		keyBytes:    []byte{5, 6, 7},
	}

	data, err := b.MarshalBinary()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(data)
	// Output:
	// [0 0 0 3 0 0 0 0 0 0 0 1 0 0 0 10 5 6 7]
}

//nolint:govet
func Example_indexEntryUnmarshalBinary() {
	var b indexEntry

	err := b.UnmarshalBinary([]byte{0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 10, 5, 6, 7})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(b)
	// Output:
	// {1 10 [5 6 7]}
}
