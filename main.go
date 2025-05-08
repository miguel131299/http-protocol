package main

import (
	"fmt"
	"io"
	"log"
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
	file, err := os.Open("messages.txt") // For read access.
	if err != nil {
		log.Fatal(err)
	}

	ch := getLinesChannel(file)

	for line := range ch {
		fmt.Printf("read: %s\n", line)
	}
}
