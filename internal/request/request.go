package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       int
}

func (r *Request) parse(data []byte) (int, error) {
	_, bytesRead, err := ParseRequestLine(data)

	if err != nil {
		return 0, err
	}
	if bytesRead == 0 {
		return 0, errors.New("Not enough data")
	}
	return 1, nil

}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	//Reades request
	request, err := io.ReadAll(reader)
	//Checks errors
	if err != nil {
		return nil, err
	}
	//Parse the data to get request line
	requestline, _, err := ParseRequestLine(request)
	//Checks errors
	if err != nil {
		return nil, err
	}
	httpRequest := &Request{
		*requestline,
		1,
	}
	return httpRequest, nil
}
func ParseRequestLine(f []byte) (*RequestLine, int, error) {
	if strings.Contains(string(f), "\r\n") == false {
		return nil, 0, nil
	}
	lines := strings.Split(string(f), "\r\n")
	requestline := strings.Split(lines[0], " ")
	if len(requestline) != 3 {
		return nil, 0, errors.New("Too many spaces in request line")
	}
	for _, i := range requestline[0] {
		if i < 'A' || i > 'Z' {
			return nil, 0, errors.New("Need Upercase Method in Request line")
		}
	}
	httpversion := strings.Split(requestline[2], "/")
	if httpversion[1] != "1.1" {
		return nil, 0, errors.New("Wrong http version")
	}
	Request := &RequestLine{
		httpversion[1],
		requestline[1],
		requestline[0],
	}
	return Request, len(lines[0]), nil
}
