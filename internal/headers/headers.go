package headers

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Headers map[string]string

const httpSepCRLF = "\r\n"
const whitespaceChars = " \n"

func isValidKey(key []byte) bool {
	if len(key) == 0 {
		return false
	}
	for _, ch := range key {
		if 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || '0' <= ch && ch <= '9' {
			continue
		}
		switch ch {
		case '!', '#', '$', '%', '&', '\'', '*', '+', '-', '.', '^', '_', '`', '|', '~':
			continue
		default:
			return false
		}
	}
	return true
}

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Get(key string) (val string, ok bool) {
	val, ok = h[strings.ToLower(key)]
	return
}
func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	if v, ok := h[key]; ok {
		h[key] = fmt.Sprintf("%v, %v", v, value)
	} else {
		h[key] = value
	}
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
	if !isValidKey(key) {
		fmt.Printf("Key: %s, len: %v\n", key, len(key))
		return 0, false, errors.New("Invalid Key")
	}
	if bytes.ContainsAny(value, whitespaceChars) {
		return 0, false, errors.New("Whitespace found inside value")
	}
	h.Set(string(key), string(value))
	return idx + len(httpSepCRLF), false, nil
}
