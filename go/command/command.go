// Binary command is a command line tool for tables.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jaeyeom/sstable/go/sstable"
)

// cat prints the list of keys and values of each path in the RecordIO.
func cat(tablePaths []string) {
	for _, tablePath := range tablePaths {
		f, err := os.Open(tablePath)
		if err != nil {
			log.Println("Error on opening path", tablePath, ":", err)
			return
		}
		defer f.Close()

		info, err := f.Stat()
		if err != nil {
			log.Println("Error on stating path", tablePath, ":", err)
			return
		}

		size := info.Size()
		if size < 0 {
			log.Println("Size is negative on path", tablePath)
			return
		}

		for c := sstable.NewRecordIOReader(f, uint64(size)); !c.Done(); c.Next() {
			fmt.Println(string(c.Entry().Key))
			fmt.Println(string(c.Entry().Value))
		}
	}
}

// lss prints the list of keys of each path in the SSTable.
func lss(tablePaths []string) {
	for _, tablePath := range tablePaths {
		f, err := os.Open(tablePath)
		if err != nil {
			log.Println("Error on opening path:", err)
		}
		defer f.Close()

		tbl, err := sstable.NewSSTable(f)
		if err != nil {
			log.Println("Error on path", tablePath, ":", err)
			return
		}

		for c := tbl.ScanFrom(nil); !c.Done(); c.Next() {
			fmt.Println(string(c.Entry().Key))
		}
	}
}

// cats prints the value of the key in the SSTable.
func cats(tablePath string, key string) {
	f, err := os.Open(tablePath)
	if err != nil {
		log.Println("Error on opening path", tablePath, ":", err)
		return
	}
	defer f.Close()

	tbl, err := sstable.NewSSTable(f)
	if err != nil {
		log.Println(err)
		return
	}

	c := tbl.ScanFrom([]byte(key))
	if c.Done() || !bytes.Equal(c.Entry().Key, []byte(key)) {
		return
	}

	value := c.Entry().Value
	fmt.Println(string(value))
}

// appendKeyValue appends the key and value to the RecordIO.
func appendKeyValue(tablePath string, key string, value string) {
	f, err := os.OpenFile(tablePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
		return
	}
	defer f.Close()

	e := sstable.Entry{Key: []byte(key), Value: []byte(value)}
	if _, err = e.WriteTo(f); err != nil {
		log.Println(err)
		return
	}
}

// convert converts a RecordIO to SSTable.
func convert(from, to string) {
	f, err := os.Open(from)
	if err != nil {
		log.Printf("Error on opening path %q: %+v", from, err)
		return
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		log.Printf("Error on stating path %q: %+v", from, err)
		return
	}

	size := info.Size()
	if size < 0 {
		log.Printf("Size is negative on path %q", from)
		return
	}

	t, err := os.Create(to)
	if err != nil {
		log.Printf("Error on creating path %q: %+v", to, err)
		return
	}

	w := sstable.NewWriter(t)
	defer w.Close()

	for c := sstable.NewRecordIOReader(f, uint64(size)); !c.Done(); c.Next() {
		if err := w.Write(*c.Entry()); err != nil {
			log.Printf("Error on writing to sstable %q: %+v", to, err)
		}
	}
}

// help prints help message. If cmd is empty, prints the list of commands.
func help(cmd string) {
	helpDetails := map[string]string{
		"cat":     "cat path [path...] - prints the all keys and values in the RecordIO",
		"lss":     "lss path [path...] - prints list of keys from each path in the SSTable",
		"cats":    "cats path key - prints the value in the SSTable",
		"append":  "append path key value - append key, value to path in RecordIO",
		"convert": "convert from to - convert a RecordIO file to an SSTable. RecordIO should be already sorted",
	}

	if cmd == "" {
		fmt.Println("Available commands are:")

		for cmd, details := range helpDetails {
			fmt.Println(cmd, ":", details)
		}
	} else {
		fmt.Println(helpDetails[cmd])
	}
}

func main() {
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 || (len(args) == 1 && args[0] == "help") {
		help("")
		return
	}

	if len(args) == 2 && args[0] == "help" {
		help(args[1])
		return
	}

	switch cmd := args[0]; cmd {
	case "cat":
		if len(args) < 2 {
			help("cat")
			return
		}

		cat(args[1:])
	case "lss":
		if len(args) < 2 {
			help("ls")
			return
		}

		lss(args[1:])
	case "cats":
		if len(args) != 3 {
			help("cat")
			return
		}

		cats(args[1], args[2])
	case "append":
		if len(args) != 4 {
			help("append")
			return
		}

		appendKeyValue(args[1], args[2], args[3])
	case "convert":
		if len(args) != 3 {
			help("cat")
			return
		}

		convert(args[1], args[2])
	}
}
