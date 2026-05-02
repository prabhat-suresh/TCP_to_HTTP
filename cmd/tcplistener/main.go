package main

import (
	"fmt"
	"httpfromtcp/internal/request"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Error: ", err)
		}
		fmt.Println("Connection Accepted")
		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatal("Error: ", err)
		}
		fmt.Printf("Request line:\n - Method: %v\n - Target: %v\n - Version: %v\n",
			req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
		fmt.Println("Connection Closed")
	}
}
