package headers

import (
	"bytes"
	"errors"
)

type Headers map[string]string

const httpSepCRLF = "\r\n"
const whitespaceChars = " \n"

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(httpSepCRLF))
	if idx == -1 {
		//CRLF not found, need to read more data before parsing
		return 0, false, nil
	}
	if idx == 0 {
		//CRLF found at beginning implies end of headers section
		return len(httpSepCRLF), true, nil
	}

	headerParts := bytes.SplitN(data[:idx], []byte(":"), 2)
	if len(headerParts) < 2 {
		return 0, false, errors.New("Malformed Header")
	}
	key, value := headerParts[0], bytes.TrimSpace(headerParts[1])
	if len(key) == 0 {
		return 0, false, errors.New("Key not found")
	}
	if bytes.ContainsAny(key, whitespaceChars) {
		return 0, false, errors.New("Whitespace found around key")
	}
	if bytes.ContainsAny(value, whitespaceChars) {
		return 0, false, errors.New("Whitespace found inside value")
	}
	h[string(key)] = string(value)
	return idx + len(httpSepCRLF), false, nil
}
