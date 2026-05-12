package main

import (
	"fmt"
	"io"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		return
	}
	defer listener.Close()
	for {
		file, err := listener.Accept()
		if err != nil {
			break
		}
		fmt.Println("Connection has been Accepted")
		Lines := getLinesChannel(file)
		// Prints lines recived
		for line := range Lines {
			fmt.Println(line)
		}
		fmt.Println("Connection closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	var Channel chan string = make(chan string)
	go func() {
		defer f.Close()
		oneLine := []byte{}
		// Reads 8 bytes at a time and fromats line by line
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
