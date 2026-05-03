package server

import (
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"log"
	"net"
	"strconv"
	"sync/atomic"
	"time"
)

type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

type Handler func(w *response.Writer, req *request.Request)

func Serve(handler Handler, port int) (*Server, error) {
	listener, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}
	server := &Server{listener: listener, handler: handler}
	go server.listen()
	return server, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	return s.listener.Close()
}
func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		fmt.Println("Connection Accepted")
		go s.handle(conn)
	}
}
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	// Require setting deadline as in my case the browser doesn't seem to send EOF
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))
	req, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Printf("Connection Closed due to invalid request: %v\n", err.Error())
		return
	}
	writer := response.Writer{Writer: conn}
	s.handler(&writer, req)
	fmt.Println("Connection Closed")
}
