package sstable

import (
	"fmt"
)

func ExampleEntryMarshalBinary() {
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

func ExampleEntryUnmarshalBinary() {
	e := Entry{}
	err := e.UnmarshalBinary([]byte{0, 0, 0, 3, 0, 0, 0, 4, 1, 2, 3, 5, 6, 7, 8})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(e)
	// Output:
	// {[1 2 3] [5 6 7 8]}
}
