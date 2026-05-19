package main

import (
	"fmt"
	"httpServer/internal/request"
	"net"
)

func main() {
	//opens the listner
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		return
	}
	//closes listner when function ends
	defer listener.Close()
	for {
		file, err := listener.Accept()
		// reads from the listner
		//checks for errors
		if err != nil {
			break
		}
		fmt.Println("Connection has been Accepted")
		//parses request
		Lines, err := request.RequestFromReader(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		// prints the request
		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", Lines.RequestLine.Method, Lines.RequestLine.RequestTarget, Lines.RequestLine.HttpVersion)
		fmt.Printf("Headers:\n")
		for x, i := range Lines.Headers {
			fmt.Printf("- %s: %s\n", x, i)
		}
		fmt.Printf("Body:\n")
		fmt.Printf("%s\n", string(Lines.Body))
		fmt.Println("Connection closed")
	}
}
