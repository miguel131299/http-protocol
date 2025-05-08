package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	raddr, _ := net.ResolveUDPAddr("udp", "localhost:42069")
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading message %v\n", err)
			break
		}

		_, err = conn.Write([]byte(msg))
		if err != nil {
			fmt.Printf("Error writing to connection: %v\n", err)
			break
		}

	}

}
