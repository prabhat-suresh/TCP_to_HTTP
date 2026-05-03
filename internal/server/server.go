package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	closed   atomic.Bool
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	server := &Server{listener: listener}
	go func() {
		server.listen()
	}()
	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}
func (s *Server) listen() {
	for {
		if s.closed.Load() {
			break
		}
		conn, err := s.listener.Accept()
		if err != nil {
			log.Fatal("Error: ", err)
		}
		fmt.Println("Connection Accepted")
		go func() {
			s.handle(conn)
		}()
	}
}
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 13\r\n\r\nHello World!\n"))
	fmt.Println("Connection Closed")
}
