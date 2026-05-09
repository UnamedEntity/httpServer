package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	file, errored := os.Open("messages.txt")
	if errored != nil {
		fmt.Println("Error Occured")
		return
	}
	defer file.Close()

	for {
		b := make([]byte, 8)
		n, err := io.ReadFull(file, b)

		if n > 0 {
			fmt.Printf("read: %s\n", b[:n])
		}
		if err == io.EOF {
			break
		}
	}
}
