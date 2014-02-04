package sstable

import (
	"fmt"
	"io/ioutil"
	"os"
)

func ExampleWriter() {
	f, _ := ioutil.TempFile("", "")
	name := f.Name()
	defer os.Remove(name)
	w := NewWriter(f)
	w.Write(Entry{[]byte{1, 2, 3}, []byte{5, 6, 7, 8}})
	w.Write(Entry{[]byte{2, 2, 3}, []byte{8, 5, 6, 7, 8}})
	w.Close()
	b, _ := ioutil.ReadFile(name)
	fmt.Println(b)
	// Output:
	// [0 0 0 2 0 0 0 1 0 0 0 0 0 0 0 47 0 0 0 3 0 0 0 4 1 2 3 5 6 7 8 0 0 0 3 0 0 0 5 2 2 3 8 5 6 7 8 0 0 0 3 0 0 0 0 0 0 0 16 0 0 0 31 1 2 3]
}
