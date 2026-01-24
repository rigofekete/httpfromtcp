package headers

import (
	"fmt"
	"strings"
	"bytes"
)

type Headers map[string]string

func NewHeaders() Headers {
	return Headers{} 
}

const crlf = "\r\n"

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, []byte(crlf))	
	if idx == -1 {
		return 0, false, nil	
	}
	if idx == 0 {
		return len(crlf), true, nil 
	}
	headerLineText := string(data[:idx])
	colonIdx := strings.Index(headerLineText, ":")
	// key needs to be lower case since the header field-name is case insensitive, as per RFC 9110 
	rawKey := strings.ToLower(headerLineText[:colonIdx])
	rawValue := headerLineText[colonIdx+1:]

	if rawKey[len(rawKey)-1] == ' ' {
		return 0, false, fmt.Errorf("malformed Header key: %s", rawKey)
	}
	key := strings.TrimSpace(rawKey)
	if !checkValidChars(key) {
		return 0, false, fmt.Errorf("invalid character/s found: %s", rawKey)
	}
	value := strings.TrimSpace(rawValue)
	if _, exists := h[key]; exists {
		h[key] = h[key] + ", " + value
	} else {
		h[key] = value
	}

	return idx + len(crlf), false, nil
}


func (h Headers) Get(key string) (string, bool) {
	value, exist := h[strings.ToLower(key)]	
	if exist {
		return value, true
	}
	fmt.Printf("Headers map key '%s' does not exist\n", key)
	return "", false
}


// helpers for field-name character validation to match RFC 9110 guidelines
var validChars = map[string]struct{}{
	"!": {},
	"#": {},
	"$": {},
	"%": {},
	"&": {},
	"'": {},
	"*": {},
	"+": {},
	"-": {},
	".": {},
	"^": {},
	"_": {},
	"`": {},
	"|": {},
	"~": {},
}

func checkValidChars(key string) bool {	
	if len(key) < 1 {
		return false
	}

	for _, c := range key {
		if (c >= 0 && c <= 9) || (c >= 'A' && c <= 'Z') ||
		(c >= 'a' && c <= 'z') {
			continue
		}
		if _, ok := validChars[string(c)]; !ok {
			return false
		} 
	}
	return true
}
