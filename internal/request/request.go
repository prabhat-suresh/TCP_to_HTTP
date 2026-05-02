package request

import (
	"bytes"
	"errors"
	"io"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type parseState int

const (
	initialized parseState = iota
	done
)

type parserState struct {
	state           parseState
	bytesToBeParsed []byte
}

const httpSepCRLF = "\r\n"
const initBufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	var req *Request
	parserState := parserState{
		state:           initialized,
		bytesToBeParsed: make([]byte, 0, initBufferSize),
	}
	buf := make([]byte, initBufferSize)
	for parserState.state != done {
		numBytesRead, err := reader.Read(buf)
		if err != nil {
			return nil, err
		}
		req, err = parserState.parse(buf[:numBytesRead])
		if err != nil {
			return nil, err
		}
	}
	return req, nil
}

func (p *parserState) parse(data []byte) (*Request, error) {
	var req Request
	p.bytesToBeParsed = append(p.bytesToBeParsed, data...)
	reqLine, numBytesParsed, err := parseRequestLine(p.bytesToBeParsed)
	if err != nil {
		return nil, err
	}
	if numBytesParsed == 0 {
		return nil, nil
	}
	req.RequestLine = *reqLine
	p.state = done
	return &req, nil
}

// TODO: convert this into a method
func parseRequestLine(plainReqLine []byte) (*RequestLine, int, error) {
	idx := bytes.Index(plainReqLine, []byte(httpSepCRLF))
	if idx == -1 {
		return nil, 0, nil
	}
	reqLine := bytes.Split(plainReqLine[:idx], []byte(" "))
	if len(reqLine) != 3 {
		return nil, 0, errors.New("Invalid number of parts in request line")
	}
	method, target, version := reqLine[0], reqLine[1], reqLine[2]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return nil, 0, errors.New("Method should consist of only Capital Letters")
		}
	}
	if string(version) != "HTTP/1.1" {
		return nil, 0, errors.New("Unsupported HTTP version")
	}
	return &RequestLine{"1.1", string(target), string(method)}, idx + len(httpSepCRLF), nil
}
