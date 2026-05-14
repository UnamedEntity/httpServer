package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	headers := strings.SplitN(string(data), "\r\n", 2)
	if len(headers) == 1 {
		return 0, false, nil
	} else if headers[0] == "" {
		return 0, true, nil
	}
	pair := strings.SplitN(headers[0], ":", 2)
	if len(pair) != 2 {
		return 0, false, errors.New("Missing \":\" ")
	} else if pair[0] != strings.TrimSpace(pair[0]) {
		return 0, false, errors.New("Too many spaces in header")
	} else {
		h[pair[0]] = strings.TrimSpace(pair[1])
		return len(headers[0]) + 2, false, nil
	}
}
