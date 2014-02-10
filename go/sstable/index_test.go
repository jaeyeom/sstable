package sstable

import (
	"bytes"
	"fmt"
	"io"
)

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

func ExampleIndexBufferWrite() {
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

func ExampleIndexEntryIndexOf() {
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

func ExampleIndexReadFrom() {
	i := index{}
	buf := bytes.NewBuffer([]byte{
		0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 234, 119, 1, 2, 3,
		0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 234, 119, 0, 0, 117, 59, 2, 3, 4,
	})
	i.ReadFrom(buf)
	fmt.Println(i)
	// Output:
	// [{0 60023 [1 2 3]} {60023 30011 [2 3 4]}]
}

func ExampleIndexReadAt() {
	i := index{}
	buf := fakeReader{
		buf: []byte{
			0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 234, 119, 1, 2, 3,
			0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 234, 119, 0, 0, 117, 59, 2, 3, 4,
		},
	}
	i.ReadAt(&buf, 0)
	fmt.Println(i)
	// Output:
	// [{0 60023 [1 2 3]} {60023 30011 [2 3 4]}]
}

func ExampleIndexWriteTo() {
	i := &index{
		{0, 60023, []byte{1, 2, 3}},
		{60023, 30011, []byte{2, 3, 4}},
	}
	buf := bytes.NewBuffer([]byte{})
	i.WriteTo(buf)
	fmt.Println(buf.Bytes())
	// Output:
	// [0 0 0 3 0 0 0 0 0 0 0 0 0 0 234 119 1 2 3 0 0 0 3 0 0 0 0 0 0 234 119 0 0 117 59 2 3 4]
}

func ExampleIndexEntryMarshalBinary() {
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

func ExampleIndexEntryUnmarshalBinary() {
	b := indexEntry{}
	err := b.UnmarshalBinary([]byte{0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 10, 5, 6, 7})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(b)
	// Output:
	// {1 10 [5 6 7]}
}
