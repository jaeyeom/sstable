package sstable

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ExampleSSTable() {
	f, _ := ioutil.TempFile("", "")
	name := f.Name()
	defer os.Remove(name)
	w := NewWriter(f)
	w.Write(Entry{[]byte{1, 2, 3}, []byte{5, 6, 7, 8}})
	w.Write(Entry{[]byte{2, 2, 3}, []byte{8, 5, 6, 7, 8}})
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
