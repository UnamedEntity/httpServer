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
	Lines := getLinesChannel(file)
	for line := range Lines {
		fmt.Printf("read: %s\n", line)
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	var Channel chan string = make(chan string)
	go func() {
		defer f.Close()
		oneLine := []byte{}
		for {
			b := make([]byte, 8)
			n, err := f.Read(b)
			for _, i := range b[:n] {
				if i == '\n' {
					Channel <- string(oneLine)
					oneLine = []byte{}
				} else {
					oneLine = append(oneLine, i)
				}
			}
			if err == io.EOF {
				if len(oneLine) > 0 {
					Channel <- string(oneLine)
				}
				break
			}
		}
		close(Channel)
	}()
	return Channel
}
