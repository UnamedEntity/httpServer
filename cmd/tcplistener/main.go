package main

import (
	"fmt"
	"httpServer/internal/request"
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
		Lines, err := request.RequestFromReader(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s", Lines.RequestLine.Method, Lines.RequestLine.RequestTarget, Lines.RequestLine.HttpVersion)
		fmt.Println("Connection closed")
	}
}
