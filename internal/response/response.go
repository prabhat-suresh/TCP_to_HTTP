package response

import (
	"fmt"
	"httpfromtcp/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	StatusOk                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

type Writer struct {
	Writer io.Writer
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	var statusLine []byte
	switch statusCode {
	case StatusOk:
		statusLine = []byte("HTTP/1.1 200 OK\r\n")
	case StatusBadRequest:
		statusLine = []byte("HTTP/1.1 400 Bad Request\r\n")
	case StatusInternalServerError:
		statusLine = []byte("HTTP/1.1 500 Internal Server Error\r\n")
	default:
		statusLine = fmt.Appendf(nil, "HTTP/1.1 %v \r\n", statusCode)
	}
	_, err := w.Writer.Write(statusLine)
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", strconv.Itoa(contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	var err error
	for k, v := range headers {
		_, err = w.Writer.Write(fmt.Appendf(nil, "%v: %v\r\n", k, v))
		if err != nil {
			return err
		}
	}
	_, err = w.Writer.Write([]byte("\r\n"))
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	return w.Writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	lineLen1, err := w.Writer.Write(fmt.Appendf(nil, "%x\r\n", len(p)))
	if err != nil {
		return 0, err
	}
	lineLen2, err := w.Writer.Write(p)
	if err != nil {
		return 0, err
	}
	w.Writer.Write([]byte("\r\n"))
	return lineLen1 + lineLen2 + len("\r\n"), nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	return w.Writer.Write([]byte("0\r\n\r\n"))
}
