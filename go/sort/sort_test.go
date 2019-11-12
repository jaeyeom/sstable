package sort

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/jaeyeom/sstable/go/sstable"
)

func ExampleSortEntries() {
	c := &SliceCursor{
		sstable.Entry{Key: []byte{3}},
		sstable.Entry{Key: []byte{2}},
		sstable.Entry{Key: []byte{4}},
		sstable.Entry{Key: []byte{1}},
	}
	f, _ := ioutil.TempFile("", "")

	name := f.Name()
	defer os.Remove(name)

	w := sstable.NewWriter(f)
	fmt.Println(SortEntries(c, 100, w))
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
	// 4 <nil>
	// Cursor is done: true
	// &{[1] []}
	// &{[2] []}
	// &{[3] []}
	// &{[4] []}
}

//nolint:funlen
func Example_sort() {
	c := &SliceCursor{
		sstable.Entry{Key: []byte{3}},
		sstable.Entry{Key: []byte{2}},
		sstable.Entry{Key: []byte{4}},
		sstable.Entry{Key: []byte{1}},
	}

	var files []*os.File
	defer func() { //nolint:wsl
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

		fmt.Println(SortEntries(c, 20, w))
		fmt.Println("Cursor is done:", c.Done())

		r, _ := os.Open(name)
		files = append(files, r)
	}

	cursors := make([]sstable.Cursor, 0, len(files))

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
	// 3 <nil>
	// Cursor is done: false
	// 1 <nil>
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

//nolint:funlen
func ExampleMerge() {
	f, _ := ioutil.TempFile("", "")

	name := f.Name()
	defer os.Remove(name)

	w := sstable.NewWriter(f)

	cs := []sstable.Cursor{&SliceCursor{
		sstable.Entry{Key: []byte{2}},
		sstable.Entry{Key: []byte{5}},
		sstable.Entry{Key: []byte{10}},
		sstable.Entry{Key: []byte{15}},
	}, &SliceCursor{
		sstable.Entry{Key: []byte{1}},
		sstable.Entry{Key: []byte{4}},
		sstable.Entry{Key: []byte{11}},
		sstable.Entry{Key: []byte{12}},
	}, &SliceCursor{
		sstable.Entry{Key: []byte{6}},
		sstable.Entry{Key: []byte{8}},
		sstable.Entry{Key: []byte{9}},
		sstable.Entry{Key: []byte{14}},
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
