package server

import (
	"fmt"
	"httpfromtcp/internal/response"
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
	response.WriteStatusLine(conn, response.StatusOk)
	response.WriteHeaders(conn, response.GetDefaultHeaders(0))
	fmt.Println("Connection Closed")
}
