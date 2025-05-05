package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("messages.txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	data := make([]byte, 8)
	var line string
	for {
		n, err := file.Read(data)

		if n > 0 {
			chunk := string(data[:n])
			parts := strings.Split(chunk, "\n")

			line += parts[0]

			for i := 1; i < len(parts); i++ {
				fmt.Printf("read: %s\n", line)
				line = parts[i]
			}
		}

		if err != nil {
			if err == io.EOF {
				if len(line) > 0 {
					// print last line
					fmt.Printf("read: %s\n", line)
				}
				break
			}

			log.Fatal(err)
		}
	}
}
