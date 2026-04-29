package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		buffer := make([]byte, 8)
		contents := make([]byte, 0)
		for {
			numBytesRead, err := f.Read(buffer)
			if err != nil {
				if errors.Is(err, io.EOF) {
					// EOF
					break
				} else {
					// unexpected extraneous error
					log.Fatal("Error: ", err)
				}
			}
			if newlineIndex := bytes.IndexByte(buffer[:numBytesRead], '\n'); newlineIndex != -1 {
				contents = append(contents, buffer[:newlineIndex]...)
				ch <- string(contents)
				contents = contents[:0]
				contents = append(contents, buffer[newlineIndex+1:numBytesRead]...)
			} else {
				contents = append(contents, buffer[:numBytesRead]...)
			}
		}
		if len(contents) > 0 {
			ch <- string(contents)
		}
	}()
	return ch
}

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	defer file.Close()
	for line := range getLinesChannel(file) {
		fmt.Printf("read: %s\n", line)
	}
}
