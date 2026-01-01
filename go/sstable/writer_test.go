package sstable

import (
	"encoding/hex"
	"fmt"
	"os"
)

func ExampleWriter() {
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
	fmt.Print(hex.Dump(b))
	// Output:
	// 00000000  00 00 00 02 00 00 00 01  00 00 00 00 00 00 00 2f  |.............../|
	// 00000010  00 00 00 03 00 00 00 04  01 02 03 05 06 07 08 00  |................|
	// 00000020  00 00 03 00 00 00 05 02  02 03 08 05 06 07 08 00  |................|
	// 00000030  00 00 03 00 00 00 00 00  00 00 10 00 00 00 1f 01  |................|
	// 00000040  02 03                                             |..|
}
