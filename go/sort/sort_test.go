package sort

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jaeyeom/sstable/go/sstable"
)

func ExampleSortEntries() {
	c := &SliceCursor{
		sstable.Entry{[]byte{3}, []byte{}},
		sstable.Entry{[]byte{2}, []byte{}},
		sstable.Entry{[]byte{4}, []byte{}},
		sstable.Entry{[]byte{1}, []byte{}},
	}
	f, _ := ioutil.TempFile("", "")
	name := f.Name()
	defer os.Remove(name)
	w := sstable.NewWriter(f)
	SortEntries(c, 100, w)
	fmt.Println("Cursor is done:", c.Done())
	outf, _ := os.Open(name)
	defer outf.Close()
	s, _ := sstable.NewSSTable(outf)
	results := s.ScanFrom([]byte{})
	if results == nil {
		fmt.Println(results)
		return
	}
	for !results.Done() {
		fmt.Println(results.Entry())
		results.Next()
	}
	// Output:
	// Cursor is done: true
	// &{[1] []}
	// &{[2] []}
	// &{[3] []}
	// &{[4] []}
}

func ExampleSort() {
	c := &SliceCursor{
		sstable.Entry{[]byte{3}, []byte{}},
		sstable.Entry{[]byte{2}, []byte{}},
		sstable.Entry{[]byte{4}, []byte{}},
		sstable.Entry{[]byte{1}, []byte{}},
	}
	var files []*os.File
	defer func() {
		for _, f := range files {
			f.Close()
		}
	}()
	for !c.Done() {
		f, _ := ioutil.TempFile("", "")
		name := f.Name()
		defer os.Remove(name)
		w := sstable.NewWriter(f)
		defer w.Close()
		SortEntries(c, 20, w)
		fmt.Println("Cursor is done:", c.Done())
		r, _ := os.Open(name)
		files = append(files, r)
	}

	var cursors []sstable.Cursor
	for _, f := range files {
		s, _ := sstable.NewSSTable(f)
		c := s.ScanFrom(nil)
		if c == nil {
			fmt.Println(c)
			return
		}
		cursors = append(cursors, c)
	}

	f, _ := ioutil.TempFile("", "")
	name := f.Name()
	defer os.Remove(name)
	w := sstable.NewWriter(f)
	if err := Merge(cursors, w); err != nil {
		fmt.Println(err)
		return
	}
	w.Close()

	f2, _ := os.Open(name)
	s, _ := sstable.NewSSTable(f2)
	results := s.ScanFrom(nil)
	if results == nil {
		fmt.Println(results)
		return
	}
	for !results.Done() {
		fmt.Println(results.Entry())
		results.Next()
	}
	// Output:
	// Cursor is done: false
	// Cursor is done: true
	// &{[1] []}
	// &{[2] []}
	// &{[3] []}
	// &{[4] []}
}

type SliceCursor []sstable.Entry

func (c *SliceCursor) Entry() *sstable.Entry {
	return &(*c)[0]
}

func (c *SliceCursor) Next() {
	*c = (*c)[1:]
}

func (c *SliceCursor) Done() bool {
	return len(*c) == 0
}

func ExampleMerge() {
	f, _ := ioutil.TempFile("", "")
	name := f.Name()
	defer os.Remove(name)
	w := sstable.NewWriter(f)
	cs := []sstable.Cursor{&SliceCursor{
		sstable.Entry{[]byte{2}, []byte{}},
		sstable.Entry{[]byte{5}, []byte{}},
		sstable.Entry{[]byte{10}, []byte{}},
		sstable.Entry{[]byte{15}, []byte{}},
	}, &SliceCursor{
		sstable.Entry{[]byte{1}, []byte{}},
		sstable.Entry{[]byte{4}, []byte{}},
		sstable.Entry{[]byte{11}, []byte{}},
		sstable.Entry{[]byte{12}, []byte{}},
	}, &SliceCursor{
		sstable.Entry{[]byte{6}, []byte{}},
		sstable.Entry{[]byte{8}, []byte{}},
		sstable.Entry{[]byte{9}, []byte{}},
		sstable.Entry{[]byte{14}, []byte{}},
	}}
	if err := Merge(cs, w); err != nil {
		fmt.Println(err)
		return
	}
	w.Close()
	f2, _ := os.Open(name)
	defer f2.Close()
	s, err := sstable.NewSSTable(f2)
	if err != nil {
		fmt.Println(err)
		return
	}
	c := s.ScanFrom([]byte{})
	if c == nil {
		fmt.Println(c)
		return
	}
	for !c.Done() {
		fmt.Println(c.Entry())
		c.Next()
	}
	// Output:
	// &{[1] []}
	// &{[2] []}
	// &{[4] []}
	// &{[5] []}
	// &{[6] []}
	// &{[8] []}
	// &{[9] []}
	// &{[10] []}
	// &{[11] []}
	// &{[12] []}
	// &{[14] []}
	// &{[15] []}
}
