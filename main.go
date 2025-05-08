package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	go func() {
		// ensure file is closed
		defer f.Close()

		data := make([]byte, 8)
		var line string
		for {
			n, err := f.Read(data)

			if n > 0 {
				chunk := string(data[:n])
				parts := strings.Split(chunk, "\n")

				line += parts[0]

				for i := 1; i < len(parts); i++ {
					ch <- line
					line = parts[i]
				}
			}

			if err != nil {
				if err == io.EOF {
					if len(line) > 0 {
						// send last line
						ch <- line
					}
					close(ch)
					break
				}
				log.Fatal(err)
			}
		}
	}()

	return ch
}

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

		ch := getLinesChannel(conn)

		for line := range ch {
			fmt.Println(line)
		}

		fmt.Println("Connection closed")
	}
}
