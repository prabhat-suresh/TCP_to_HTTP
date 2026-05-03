package server

import (
	"bytes"
	"fmt"
	"httpfromtcp/internal/request"
	"httpfromtcp/internal/response"
	"io"
	"log"
	"net"
	"strconv"
	"sync/atomic"
)

type Server struct {
	listener net.Listener
	handler  Handler
	closed   atomic.Bool
}

type HandlerError struct {
	StatusCode response.StatusCode
	Msg        string
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

func (h *HandlerError) writeHandlerErrorTo(w io.Writer) {
	w.Write(fmt.Appendf(nil, "HTTP/1.1 %v %v\r\n", h.StatusCode, h.Msg))
}

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
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError{
			StatusCode: response.StatusBadRequest,
			Msg:        err.Error(),
		}
		hErr.writeHandlerErrorTo(conn)
		return
	}
	handlerBuf := bytes.NewBuffer([]byte{})
	handlerError := s.handler(handlerBuf, req)
	if handlerError != nil {
		handlerError.writeHandlerErrorTo(conn)
		fmt.Println("Connection Closed")
		return
	}
	response.WriteStatusLine(conn, response.StatusOk)
	response.WriteHeaders(conn, response.GetDefaultHeaders(handlerBuf.Len()))
	conn.Write(handlerBuf.Bytes())
	fmt.Println("Connection Closed")
}
