package shard

import (
	"crypto/sha512"
	"fmt"
	"io"
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
	defer func() { _ = os.RemoveAll(name) }() //nolint:wsl

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

	_ = w.Close() // Changed this line

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

func ExamplePrefixSum64Hash_Sum64() {
	hasher := &PrefixSum64Hash{sha512.New()}
	hasher.Write([]byte("hello"))
	fmt.Println(hasher.Sum64())
	// Output: 11200964803485168504
}

func ExampleNewOSFileWriterFactory() {
	tempDir, err := ioutil.TempDir("", "examplefactory")
	if err != nil {
		fmt.Println("Failed to create temp dir:", err)
		return
	}
	defer func() { _ = os.RemoveAll(tempDir) }()

	factoryPrefix := "testfactory-"
	factory := NewOSFileWriterFactory(path.Join(tempDir, factoryPrefix))

	numFiles := 3
	for i := 0; i < numFiles; i++ {
		writer := factory(i, numFiles)
		_, err := fmt.Fprintf(writer, "data for file %d", i)
		if err != nil {
			fmt.Printf("Error writing to file %d: %v\n", i, err)
		}
		if closer, ok := writer.(io.Closer); ok {
			err := closer.Close()
			if err != nil {
				fmt.Printf("Error closing file %d: %v\n", i, err)
			}
		}
	}

	for i := 0; i < numFiles; i++ {
		fileName := fmt.Sprintf("%s%05d-of-%05d", factoryPrefix, i, numFiles)
		filePath := path.Join(tempDir, fileName)
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			fmt.Printf("Error reading file %s: %v\n", fileName, err)
			continue
		}
		fmt.Printf("%s:%s\n", fileName, string(content))
	}

	// Output:
	// testfactory-00000-of-00003:data for file 0
	// testfactory-00001-of-00003:data for file 1
	// testfactory-00002-of-00003:data for file 2
}
