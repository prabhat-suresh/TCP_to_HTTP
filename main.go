package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	file, err := os.Open("messages.txt")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	buffer := make([]byte, 8)
	for {
		numBytesRead, err := file.Read(buffer)
		if err != nil {
			if errors.Is(err, io.EOF) {
				// EOF
				break
			} else {
				// unexpected extraneous error
				log.Fatal("Error: ", err)
			}
		}
		fmt.Printf("read: %s\n", buffer[:numBytesRead])
	}
}
