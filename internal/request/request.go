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

const (
	stateinit = iota
	statedone
)

func (r *Request) parse(data []byte) (int, error) {
	if r.state == stateinit {
		requestLine, byteRead, err := ParseRequestLine(data)
		if err != nil {
			return 0, err
		}
		if byteRead == 0 {
			return 0, nil
		}
		r.RequestLine = *requestLine
		r.state = statedone
		return byteRead, nil
	}
	if r.state == statedone {
		return 0, errors.New("error: trying to read data in a done state")
	}
	return 0, errors.New("unknown state")
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	httpRequest := Request{}
	parseTo := 0
	chunck := make([]byte, 8)
	for {
		if parseTo == len(chunck) {
			newbuff := make([]byte, len(chunck)*2)
			copy(newbuff, chunck)
			chunck = newbuff
		}
		//Reades request
		request, err := reader.Read(chunck[parseTo:])
		parseTo += request
		bytes, errors := httpRequest.parse(chunck[:parseTo])
		copy(chunck, chunck[bytes:parseTo])
		parseTo -= bytes
		bytes = 0
		//Checks errors
		if errors != nil {
			return nil, errors
		}
		if httpRequest.state == 1 {
			break
		}
		if err == io.EOF {
			break
		}
	}
	Request := &httpRequest
	return Request, nil
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
