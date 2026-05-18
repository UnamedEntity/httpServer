package request

import (
	"errors"
	"httpServer/internal/headers"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
	state       requestState
	Headers     headers.Headers
	Body        []byte
}

type requestState int

const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateDone
)

func (r *Request) parse(data []byte) (int, error) {
	totalByets := 0
	if r.state == requestStateInitialized {
		requestLine, byteRead, err := ParseRequestLine(data)
		// checks errors
		if err != nil {
			return 0, err
		}
		if byteRead == 0 {
			return 0, nil
		}
		// assigns the request line to the struct
		r.RequestLine = *requestLine
		r.state = requestStateParsingHeaders
		return byteRead, nil
	}

	//Parse headers
	if r.state == requestStateParsingHeaders {
		for r.state == requestStateParsingHeaders {
			n, done, err := r.Headers.Parse(data[totalByets:])
			if err != nil {
				return 0, err
			}
			if n == 0 {
				break
			}
			if done == true {
				r.state = requestStateDone
			}
			totalByets += n

		}
		return totalByets, nil
	}

	// error
	if r.state == requestStateDone {
		return 0, errors.New("error: trying to read data in a done state")
	} else {
		return 0, errors.New("unknown state")
	}

}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	httpRequest := Request{
		state:   requestStateInitialized,
		Headers: headers.NewHeaders(),
	}

	parseTo := 0
	chunck := make([]byte, 8)
	for {
		//reallocate array when max reached
		if parseTo == len(chunck) {
			newbuff := make([]byte, len(chunck)*2)
			copy(newbuff, chunck)
			chunck = newbuff
		}
		//Reades request
		request, err := reader.Read(chunck[parseTo:])
		parseTo += request
		bytes, errored := httpRequest.parse(chunck[:parseTo])
		copy(chunck, chunck[bytes:parseTo])
		parseTo -= bytes
		bytes = 0
		//Checks errors
		if errored != nil {
			return nil, errored
		}
		if httpRequest.state == requestStateDone {
			break
		}
		if err == io.EOF {
			if parseTo == 0 {
				return nil, errors.New("Empty headers")
			} else {
				break
			}
		}
	}
	// create a pointer to request
	Request := &httpRequest
	//checks for empty headers
	if len(Request.Headers) == 0 {
		return nil, errors.New("Empty headers")
	}
	return Request, nil
}
func ParseRequestLine(f []byte) (*RequestLine, int, error) {
	// check to see if it contains the CRLF
	if strings.Contains(string(f), "\r\n") == false {
		return nil, 0, nil
	}
	// splits into the request line
	lines := strings.Split(string(f), "\r\n")
	// splits into version,host, method
	requestline := strings.Split(lines[0], " ")
	//checks fo errors
	if len(requestline) != 3 {
		return nil, 0, errors.New("Too many spaces in request line")
	}
	for _, i := range requestline[0] {
		if i < 'A' || i > 'Z' {
			return nil, 0, errors.New("Need Upercase Method in Request line")
		}
	}
	//splits to find version
	httpversion := strings.Split(requestline[2], "/")
	//checks version
	if httpversion[1] != "1.1" {
		return nil, 0, errors.New("Wrong http version")
	}
	// creates a pointer to a request line
	Request := &RequestLine{
		httpversion[1],
		requestline[1],
		requestline[0],
	}
	return Request, len(lines[0]) + 2, nil
}
