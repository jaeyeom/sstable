package sstable

import (
	"bytes"
	"fmt"
)

func ExampleNewRecordIOReader() {
	// 1. Prepare entries
	entry1 := Entry{Key: []byte("keyA"), Value: []byte("valueA")}
	entry2 := Entry{Key: []byte("keyB"), Value: []byte("valueBB")} // Different length value

	// 2. Marshal entries into a buffer
	var dataBuf bytes.Buffer
	data1Bytes, err := entry1.MarshalBinary()
	if err != nil {
		fmt.Println("Error marshalling entry1:", err)
		return
	}
	dataBuf.Write(data1Bytes)

	data2Bytes, err := entry2.MarshalBinary()
	if err != nil {
		fmt.Println("Error marshalling entry2:", err)
		return
	}
	dataBuf.Write(data2Bytes)

	serializedEntries := dataBuf.Bytes()

	// 3. Create a bytes.NewReader (which implements io.Reader and io.ReaderAt)
	reader := bytes.NewReader(serializedEntries)

	// 4. Calculate the total size of the serialized data
	totalSize := uint64(len(serializedEntries))

	// 5. Call NewRecordIOReader
	// NewRecordIOReader takes an io.Reader, but CursorToOffset (which it returns)
	// will try to use it as io.ReaderAt if possible. bytes.NewReader supports this.
	cursor := NewRecordIOReader(reader, totalSize)

	// 6. Loop through entries using the cursor
	for !cursor.Done() {
		currentEntry := cursor.Entry()
		if currentEntry == nil {
			// Should not happen in this example if data is valid and size is correct
			fmt.Println("Error: current entry is nil, but not done.")
			break
		}
		fmt.Printf("Key: %s, Value: %s\n", string(currentEntry.Key), string(currentEntry.Value))
		cursor.Next()
	}

	// Output:
	// Key: keyA, Value: valueA
	// Key: keyB, Value: valueBB
}
