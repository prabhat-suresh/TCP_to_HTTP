package request

import (
	"errors"
	"io"
	"log"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

const sep = "\r\n"

func RequestFromReader(reader io.Reader) (*Request, error) {
	var req Request
	plainReq, err := io.ReadAll(reader)
	if err != nil {
		log.Fatal(err)
	}
	plainReqLines := strings.Split(string(plainReq), sep)
	reqLine, _, err := parseRequestLine(plainReqLines)
	if err != nil {
		return nil, err
	}
	req.RequestLine = *reqLine
	return &req, nil
}

func parseRequestLine(plainReqLines []string) (*RequestLine, []string, error) {
	reqLine := strings.Split(plainReqLines[0], " ")
	if len(reqLine) != 3 {
		return nil, plainReqLines, errors.New("Invalid number of parts in request line")
	}
	method, target, version := reqLine[0], reqLine[1], reqLine[2]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, plainReqLines, errors.New("Method should consist of only Capital Letters")
		}
	}
	if version != "HTTP/1.1" {
		return nil, plainReqLines, errors.New("Unsupported HTTP version")
	}
	return &RequestLine{"1.1", target, method}, plainReqLines[1:], nil
}
