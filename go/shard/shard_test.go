package shard

import (
	"crypto/sha512"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

func ExampleWriter() {
	name, err := ioutil.TempDir("", "test")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer os.RemoveAll(name) //nolint:wsl

	w := NewWriter(5, &PrefixSum64Hash{sha512.New()}, NewOSFileWriterFactory(path.Join(name, "test-")))
	records := []string{
		"test0",
		"test1",
		"test2",
		"test3",
		"test4",
		"test5",
		"test6",
		"test7",
		"test1",
	}

	for _, rec := range records {
		if _, err := w.Write([]byte(rec)); err != nil {
			fmt.Println(err)
		}
	}

	w.Close()

	for i := 0; i < 5; i++ {
		filename := fmt.Sprintf("test-%05d-of-00005", i)
		b, _ := ioutil.ReadFile(path.Join(name, filename))
		fmt.Printf("%s:%s\n", filename, string(b))
	}
	// Output:
	// test-00000-of-00005:test0test4
	// test-00001-of-00005:test7
	// test-00002-of-00005:test5test6
	// test-00003-of-00005:
	// test-00004-of-00005:test1test2test3test1
}
