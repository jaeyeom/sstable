package sstable

import (
	"bytes"
	"fmt"
	// "io" // Will add if ReadEntry or other functions require it directly for type matching.
	// For now, bytes.Reader and the CursorToOffset itself will handle reader interfaces.
)

func ExampleCursorToOffset() {
	// 1. Prepare entries
	entry1 := Entry{Key: []byte("key1"), Value: []byte("value1")}
	entry2 := Entry{Key: []byte("key2"), Value: []byte("value22")} // Different length value

	// 2. Marshal entries into a buffer
	var buf bytes.Buffer
	data1, err := entry1.MarshalBinary()
	if err != nil {
		fmt.Println("Error marshalling entry1:", err)
		return
	}
	buf.Write(data1)

	data2, err := entry2.MarshalBinary()
	if err != nil {
		fmt.Println("Error marshalling entry2:", err)
		return
	}
	buf.Write(data2)

	serializedData := buf.Bytes()

	// 3. Calculate endOffset (total size of serialized data)
	// This could also be entry1.Size() + entry2.Size()
	endOffset := uint64(len(serializedData))

	// 4. Create a bytes.NewReader
	reader := bytes.NewReader(serializedData)

	// 5. Create CursorToOffset instance
	// Reader can be io.Reader or io.ReaderAt. bytes.NewReader implements both.
	// For this example, let's ensure it's treated as a general io.Reader
	// by the cursor, though CursorToOffset will detect it as io.ReaderAt if not type asserted.
	// The implementation of CursorToOffset.Entry() prefers io.ReaderAt if available.
	cursor := &CursorToOffset{
		reader:    reader, // bytes.NewReader is also an io.ReaderAt
		offset:    0,
		endOffset: endOffset,
		entry:     nil, // Starts with no entry loaded
	}

	// 6. Loop through entries
	for !cursor.Done() {
		currentEntry := cursor.Entry()
		if currentEntry == nil {
			// This might happen if ReadEntry/ReadEntryAt fails,
			// or if Done() condition is met but loop condition was already checked.
			// Or if endOffset is 0.
			// Given the Done() logic, if Entry() returns nil, Done() should usually be true.
			// Let's assume valid entries for example purposes.
			fmt.Println("Error: current entry is nil, but not done.")
			break
		}
		fmt.Printf("Key: %s, Value: %s\n", string(currentEntry.Key), string(currentEntry.Value))
		cursor.Next()
	}

	// Output:
	// Key: key1, Value: value1
	// Key: key2, Value: value22
}
