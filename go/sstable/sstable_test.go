package sstable

import (
	"bytes"
	"fmt"
	"os"
)

func ExampleSSTable() {
	f, _ := os.CreateTemp("", "")

	name := f.Name()
	defer os.Remove(name)

	w := NewWriter(f)

	entries := []Entry{
		{Key: []byte{1, 2, 3}, Value: []byte{5, 6, 7, 8}},
		{Key: []byte{2, 2, 3}, Value: []byte{8, 5, 6, 7, 8}},
	}
	for _, entry := range entries {
		if err := w.Write(entry); err != nil {
			fmt.Println(err)
		}
	}

	w.Close()

	f2, _ := os.Open(name)
	defer f2.Close()

	s, _ := NewSSTable(f2)

	c := s.ScanFrom([]byte{1, 2, 3})
	if c == nil {
		fmt.Println(c)
		return
	}

	for !c.Done() {
		fmt.Println(c.Entry())
		c.Next()
	}

	fmt.Println("---")

	c = s.ScanFrom([]byte{1, 2, 3, 0})
	for !c.Done() {
		fmt.Println(c.Entry())
		c.Next()
	}
	// Output:
	// &{[1 2 3] [5 6 7 8]}
	// &{[2 2 3] [8 5 6 7 8]}
	// ---
	// &{[2 2 3] [8 5 6 7 8]}
}

func ExampleSSTable_reader() {
	f, _ := os.CreateTemp("", "")

	name := f.Name()
	defer os.Remove(name)

	w := NewWriter(f)

	entries := []Entry{
		{Key: []byte{1, 2, 3}, Value: []byte{5, 6, 7, 8}},
		{Key: []byte{2, 2, 3}, Value: []byte{8, 5, 6, 7, 8}},
	}
	for _, entry := range entries {
		if err := w.Write(entry); err != nil {
			fmt.Println(err)
		}
	}

	w.Close()

	b, _ := os.ReadFile(name)
	// bytes.Buffer does not support random access.
	s, _ := NewSSTable(bytes.NewBuffer(b))

	c := s.ScanFrom([]byte{1, 2, 3})
	if c == nil {
		fmt.Println(c)
		return
	}

	for !c.Done() {
		fmt.Println(c.Entry())
		c.Next()
	}
	// Output:
	// &{[1 2 3] [5 6 7 8]}
	// &{[2 2 3] [8 5 6 7 8]}
}
