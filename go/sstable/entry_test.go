package sstable

import (
	"fmt"
	"io"
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

type fakeReader struct {
	buf    []byte
	offset int
}

func (f *fakeReader) Read(p []byte) (n int, err error) {
	if f.offset >= len(f.buf) {
		return 0, io.EOF
	}
	n = copy(p, f.buf[f.offset:])
	f.offset += n
	return n, nil
}

func (f *fakeReader) ReadAt(b []byte, offset int64) (n int, err error) {
	if int64(len(f.buf)) < offset {
		return 0, io.EOF
	}
	if int64(len(f.buf)) <= offset+int64(len(b)) {
		return copy(b, f.buf[offset:]), io.EOF
	}
	return copy(b, f.buf[offset:]), nil
}

func ExampleReadEntry() {
	f := fakeReader{buf: []byte{0, 0, 0, 3, 0, 0, 0, 4, 1, 2, 3, 5, 6, 7, 8}}
	e, _ := readEntry(&f)
	fmt.Println(e)
	// Output:
	// &{[1 2 3] [5 6 7 8]}
}

func ExampleReadEntryAt() {
	f := fakeReader{buf: []byte{0, 0, 0, 3, 0, 0, 0, 4, 1, 2, 3, 5, 6, 7, 8}}
	e, _ := readEntryAt(&f, 0)
	fmt.Println(e)
	// Output:
	// &{[1 2 3] [5 6 7 8]}
}
