package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"net"
	"os"
)

func main() {
	port := ":42069"
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("Error starting server:", err)
		os.Exit(1)
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		fmt.Printf("Connection accepted on port %s\n", port)

		r, err := request.RequestFromReader(conn)
		if err != nil {
			fmt.Println("Error reading request:", err)
		}

		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)

		fmt.Println("Connection closed")
	}
}
