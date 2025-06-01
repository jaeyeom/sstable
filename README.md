# SSTable

An SSTable (Sorted String Table) is a persistent, ordered, immutable map from keys to values, where both keys and values are arbitrary byte strings. SSTables are a fundamental component in many large-scale distributed storage systems, such as Bigtable, Cassandra, and LevelDB.

They are designed for efficient storage and retrieval of large amounts of data. Data in an SSTable is sorted by key, which allows for fast lookups and range scans. SSTables are immutable, meaning once written, they cannot be modified. This simplifies concurrency control and caching. Updates and deletions are typically handled by creating new SSTables and merging them with existing ones during a compaction process.

## Features of this Implementation

This Go library provides an implementation of SSTables with the following features:

*   **Compatibility:** Fully compatible with the SSTable format defined by [https://github.com/mariusaeriksen/sstable](https://github.com/mariusaeriksen/sstable).
*   **Sorted Key-Value Storage:** Stores key-value pairs sorted by key, allowing for efficient lookups and range scans.
*   **Immutable Files:** Once an SSTable is written, it is immutable. This simplifies data management and concurrency.
*   **Data Blocks and Index:** Data is stored in blocks, and an index is created to quickly locate the block containing a specific key.
*   **Header:** Each SSTable file starts with a header containing metadata like version, number of blocks, and the offset of the index.
*   **Sequential Writes:** Optimized for sequential write patterns when creating SSTables.
*   **Random Access Reads:** Supports efficient random access to data through key lookups and range scans.
*   **Cursor for Iteration:** Provides a cursor mechanism to iterate over key-value pairs within a specified range.
*   **Flexible I/O:** Works with `io.Reader`, `io.ReaderAt`, `io.Writer`, `io.WriteSeeker`, and `io.WriterAt` interfaces, allowing for flexibility in how data is read from and written to storage.

## Usage Examples

### Creating an SSTable

```go
package main

import (
	"log"
	"os"

	"github.com/your-username/sstable/go/sstable" // Assuming this is the import path
)

func main() {
	// Create a new file for the SSTable
	f, err := os.Create("mydata.sstable")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a new SSTable writer
	writer := sstable.NewWriter(f)

	// Define some key-value pairs (must be sorted by key)
	entries := []sstable.Entry{
		{Key: []byte("apple"), Value: []byte("A fruit that grows on trees")},
		{Key: []byte("banana"), Value: []byte("A long curved fruit")},
		{Key: []byte("cherry"), Value: []byte("A small, round, red or black fruit")},
	}

	// Write entries to the SSTable
	for _, entry := range entries {
		if err := writer.Write(entry); err != nil {
			log.Fatalf("Failed to write entry: %v", err)
		}
	}

	// Close the writer to finalize the SSTable (writes the index and header)
	if err := writer.Close(); err != nil {
		log.Fatalf("Failed to close writer: %v", err)
	}

	log.Println("SSTable created successfully: mydata.sstable")
}
```

### Reading from an SSTable

```go
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/your-username/sstable/go/sstable" // Assuming this is the import path
)

func main() {
	// Open an existing SSTable file
	f, err := os.Open("mydata.sstable")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Create a new SSTable reader
	reader, err := sstable.NewSSTable(f)
	if err != nil {
		log.Fatalf("Failed to create SSTable reader: %v", err)
	}

	// Example 1: Scan all entries
	fmt.Println("Scanning all entries:")
	scannerAll := reader.ScanFrom(nil) // Start scan from the beginning
	for !scannerAll.Done() {
		entry := scannerAll.Entry()
		fmt.Printf("  Key: %s, Value: %s\n", string(entry.Key), string(entry.Value))
		scannerAll.Next()
	}
	if err := scannerAll.Error(); err != nil {
		log.Fatalf("Error during scan: %v", err)
	}

	// Example 2: Scan from a specific key
	fmt.Println("\nScanning from 'banana':")
	scannerFrom := reader.ScanFrom([]byte("banana"))
	for !scannerFrom.Done() {
		entry := scannerFrom.Entry()
		fmt.Printf("  Key: %s, Value: %s\n", string(entry.Key), string(entry.Value))
		scannerFrom.Next()
	}
	if err := scannerFrom.Error(); err != nil {
		log.Fatalf("Error during scan from key: %v", err)
	}

	// Note: For specific key lookups, you would typically iterate
	// with ScanFrom and stop when the desired key is found or passed.
	// A direct Get(key) method is not explicitly shown in the provided sstable.go,
	// but ScanFrom can be used to achieve this.
}
```

## Build and Test

This project uses Go modules.

### Building

To build the library and command-line utilities (if any):

```bash
go build ./...
```

### Testing

To run the unit tests:

```bash
go test ./...
```

You can also run tests with coverage:

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

This project includes a GitHub Actions workflow (in `.github/workflows/go.yml`) that automatically builds and tests the code on pushes and pull requests.

## Contributing

Contributions are welcome! If you find a bug, have a feature request, or want to contribute code, please follow these steps:

1.  **Check for existing issues:** Before opening a new issue, please search the existing issues to see if your problem or idea has already been discussed.
2.  **Open an issue:** If you don't find an existing issue, open a new one. Provide a clear description of the bug or feature request.
3.  **Fork the repository:** Create your own fork of the repository on GitHub.
4.  **Create a branch:** Create a new branch in your fork for your changes. Use a descriptive branch name (e.g., `fix-scan-bug`, `add-compression-feature`).
5.  **Make your changes:** Implement your changes, ensuring that the code adheres to the project's style and that all tests pass.
6.  **Add tests:** If you're adding a new feature or fixing a bug, please include unit tests that cover your changes.
7.  **Commit your changes:** Make clear and concise commit messages.
8.  **Push your changes:** Push your branch to your fork on GitHub.
9.  **Submit a pull request:** Open a pull request from your branch to the main repository. Provide a clear description of your changes in the pull request.

We will review your pull request and provide feedback as soon as possible.

## License

This project is licensed under the [MIT License](LICENSE).
