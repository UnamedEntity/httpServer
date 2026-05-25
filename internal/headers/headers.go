package headers

import (
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Get(key string) (value string, err error) {
	value, ok := h[strings.ToLower(key)]
	//checks if key exist
	if ok == false {
		return "", errors.New("No content length")
	} else {
		return value, nil
	}

}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	headers := strings.SplitN(string(data), "\r\n", 2)
	validCharacters := "!#$%&'*+-.^_`|~ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	// Check if there is a header
	if len(headers) == 1 {
		return 0, false, nil
	} else if headers[0] == "" {
		return 0, true, nil
	}
	//Splits into key value pair
	pair := strings.SplitN(headers[0], ":", 2)
	if len(pair) != 2 {
		return 0, false, errors.New("Missing \":\" ")
	}
	//Checks for vaild characters in the header
	for _, i := range pair[0] {
		if strings.Contains(validCharacters, string(i)) == false {
			return 0, false, errors.New("Invalid Characters")
		}
	}
	// Increments the value if it exists and creates a new key value pair if not
	if _, ok := h[strings.ToLower(pair[0])]; ok == false {
		h[strings.ToLower(pair[0])] = strings.TrimSpace(pair[1])
	} else {
		h[strings.ToLower(pair[0])] += ", " + strings.TrimSpace(pair[1])
	}
	//returns bytes consumed
	return len(headers[0]) + 2, false, nil
}
