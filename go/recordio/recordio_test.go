package recordio

import (
	"fmt"
)

type FakeWriteCloser struct{}

func (f FakeWriteCloser) Write(data []byte) (int, error) {
	return len(data), nil
}

func (f FakeWriteCloser) Close() error {
	fmt.Println("Closed")
	return nil
}

func ExampleWriteCloser() {
	wc := NewWriteCloser(&FakeWriteCloser{})
	wc.Close()
	// Output:
	// Closed
}
