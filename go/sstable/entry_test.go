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

func ExampleEntry_WriteTo() {
	e := Entry{
		Key:   []byte{1, 2},       // Key length 2
		Value: []byte{3, 4, 5}, // Value length 3
	}

	var buf bytes.Buffer
	_, err := e.WriteTo(&buf)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Optionally print n, but the main thing is the buffer content
	// fmt.Printf("Bytes written: %d\n", n)
	// For example consistency, usually only the primary output is shown.

	fmt.Println(buf.Bytes())
	// Expected format: [0 0 0 keyLen 0 0 0 valLen keyBytes valBytes]
	// KeyLen = 2 -> [0 0 0 2]
	// ValLen = 3 -> [0 0 0 3]
	// Key    = [1 2]
	// Value  = [3 4 5]
	// Output:
	// [0 0 0 2 0 0 0 3 1 2 3 4 5]
}

func ExampleEntry_Size() {
	e := Entry{
		Key:   []byte("test"),  // length 4
		Value: []byte("example"), // length 7
	}

	size := e.Size()
	fmt.Println(size)
	// Expected: 8 (for length fields) + 4 (len("test")) + 7 (len("example")) = 19
	// Output:
	// 19
}

func ExampleReadEntryAt() {
	f := bytes.NewReader([]byte{0, 0, 0, 3, 0, 0, 0, 4, 1, 2, 3, 5, 6, 7, 8})
	e, _ := ReadEntryAt(f, 0)
	fmt.Println(e)
	// Output:
	// &{[1 2 3] [5 6 7 8]}
}
