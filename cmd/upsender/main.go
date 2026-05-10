package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	address, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		return
	}
	connection, errored := net.DialUDP("udp", nil, address)
	if errored != nil {
		return
	}
	defer connection.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println(">")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		_, errored := connection.Write([]byte(line))
		if errored != nil {
			fmt.Println(errored)
			break
		}

	}
}
