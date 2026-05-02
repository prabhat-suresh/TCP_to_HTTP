package request

import (
	"bytes"
	"errors"
	"httpfromtcp/internal/headers"
	"io"
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type parseState int

const (
	initialized parseState = iota
	requestStateParsingHeaders
	requestStateDone
)

type parserState struct {
	state           parseState
	bytesToBeParsed []byte
}

const httpSepCRLF = "\r\n"
const initBufferSize = 8

func RequestFromReader(reader io.Reader) (*Request, error) {
	var req Request
	req.Headers = headers.NewHeaders()
	parserState := parserState{
		state:           initialized,
		bytesToBeParsed: make([]byte, 0, initBufferSize),
	}
	buf := make([]byte, initBufferSize)
	for parserState.state != requestStateDone {
		numBytesRead, err := reader.Read(buf)
		if err != nil {
			return nil, err
		}
		err = parserState.parse(&req, buf[:numBytesRead])
		if err != nil {
			return nil, err
		}
	}
	return &req, nil
}

func (p *parserState) parse(req *Request, data []byte) error {
	// When append results in allocation of new underlying array
	// due to insufficient capacity, it will only copy the data still
	// being referenced by the slice. Hence I'm not copying back data
	// to be parsed to the start of the buffer each time. Instead I
	// just update the slice to reference the part which needs to be
	// parsed. Waiting for a reallocation to clean up the unused parts
	p.bytesToBeParsed = append(p.bytesToBeParsed, data...)
	switch p.state {
	case initialized:
		numBytesParsed, err := req.RequestLine.parseRequestLine(p.bytesToBeParsed)
		if err != nil {
			return err
		}
		if numBytesParsed == 0 {
			return nil
		}
		p.state = requestStateParsingHeaders
		p.bytesToBeParsed = p.bytesToBeParsed[numBytesParsed:]
	case requestStateParsingHeaders:
		numBytesParsed, done, err := req.Headers.Parse(p.bytesToBeParsed)
		if err != nil {
			return err
		}
		if numBytesParsed == 0 {
			return nil
		}
		p.bytesToBeParsed = p.bytesToBeParsed[numBytesParsed:]
		if done {
			p.state = requestStateDone
		}
	}
	return nil
}

func (r *RequestLine) parseRequestLine(plainReqLine []byte) (int, error) {
	idx := bytes.Index(plainReqLine, []byte(httpSepCRLF))
	if idx == -1 {
		return 0, nil
	}
	reqLine := bytes.Split(plainReqLine[:idx], []byte(" "))
	if len(reqLine) != 3 {
		return 0, errors.New("Invalid number of parts in request line")
	}
	method, target, version := reqLine[0], reqLine[1], reqLine[2]
	for _, c := range method {
		if c < 'A' || c > 'Z' {
			return 0, errors.New("Method should consist of only Capital Letters")
		}
	}
	if string(version) != "HTTP/1.1" {
		return 0, errors.New("Unsupported HTTP version")
	}
	r.Method, r.RequestTarget, r.HttpVersion = string(method), string(target), "1.1"
	return idx + len(httpSepCRLF), nil
}
