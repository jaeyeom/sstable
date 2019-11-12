package sstable

import (
	"bytes"
	"fmt"
)

func ExampleEntry_MarshalBinary() {
	e := Entry{
		Key:   []byte{1, 2, 3},
		Value: []byte{5, 6, 7, 8},
	}

	data, err := e.MarshalBinary()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(data)
	// Output:
	// [0 0 0 3 0 0 0 4 1 2 3 5 6 7 8]
}

func ExampleEntry_UnmarshalBinary() {
	var e Entry

	err := e.UnmarshalBinary([]byte{0, 0, 0, 3, 0, 0, 0, 4, 1, 2, 3, 5, 6, 7, 8})
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(e)
	// Output:
	// {[1 2 3] [5 6 7 8]}
}

func ExampleReadEntry() {
	f := bytes.NewReader([]byte{0, 0, 0, 3, 0, 0, 0, 4, 1, 2, 3, 5, 6, 7, 8})
	e, _ := ReadEntry(f)
	fmt.Println(e)
	// Output:
	// &{[1 2 3] [5 6 7 8]}
}

func ExampleReadEntryAt() {
	f := bytes.NewReader([]byte{0, 0, 0, 3, 0, 0, 0, 4, 1, 2, 3, 5, 6, 7, 8})
	e, _ := ReadEntryAt(f, 0)
	fmt.Println(e)
	// Output:
	// &{[1 2 3] [5 6 7 8]}
}
