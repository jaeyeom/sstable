package shard

import (
	"crypto/sha1"
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
	defer os.RemoveAll(name)
	w := NewWriter(5, &PrefixSum64Hash{sha1.New()}, NewOSFileWriterFactory(path.Join(name, "test-")))
	w.Write([]byte("test0"))
	w.Write([]byte("test1"))
	w.Write([]byte("test2"))
	w.Write([]byte("test3"))
	w.Write([]byte("test4"))
	w.Write([]byte("test1"))
	w.Close()
	for i := 0; i < 5; i++ {
		filename := fmt.Sprintf("test-%05d-of-00005", i)
		b, _ := ioutil.ReadFile(path.Join(name, filename))
		fmt.Printf("%s:%s\n", filename, string(b))
	}
	// Output:
	// test-00000-of-00005:test1test1
	// test-00001-of-00005:
	// test-00002-of-00005:test0test3test4
	// test-00003-of-00005:
	// test-00004-of-00005:test2
}
