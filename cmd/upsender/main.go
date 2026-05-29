package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	// gets adrress
	address, err := net.ResolveUDPAddr("udp", "localhost:42069")
	//checks for erros
	if err != nil {
		return
	}
	//opens connection
	connection, errored := net.DialUDP("udp", nil, address)
	// checks for errors
	if errored != nil {
		return
	}
	// close connection when functions ends
	defer connection.Close()
	// creates a reader
	reader := bufio.NewReader(os.Stdin)
	for {
		// Reads until \n
		fmt.Println(">")
		line, err := reader.ReadString('\n')
		// checks for errors
		if err != nil {
			fmt.Println(err)
			break
		}
		// write the line to connection
		_, errored := connection.Write([]byte(line))
		//checks for errors
		if errored != nil {
			fmt.Println(errored)
			break
		}

	}
}
