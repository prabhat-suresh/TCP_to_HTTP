package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", "localhost:42069")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	stdinReader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("> ")
		line, err := stdinReader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
		}
		_, err = conn.Write([]byte(line))
		if err != nil {
			fmt.Println(err)
		}
	}
}
