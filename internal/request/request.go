package request

import (
	"errors"
	"httpServer/internal/headers"
	"io"
	"strconv"
	"strings"
)

// request struct
type Request struct {
	RequestLine RequestLine
	state       requestState
	Headers     headers.Headers
	Body        []byte
}

type requestState int

// enum
const (
	requestStateInitialized requestState = iota
	requestStateParsingHeaders
	requestStateParsingBody
	requestStateDone
)

func (r *Request) parse(data []byte) (int, error) {
	// counts total bytes
	totalByets := 0
	//state machine
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
		//changes state for next parse
		r.state = requestStateParsingHeaders
		return byteRead, nil
	}

	// Parse headers
	if r.state == requestStateParsingHeaders {
		//loops until state change
		for r.state == requestStateParsingHeaders {
			// parse headers
			n, done, err := r.Headers.Parse(data[totalByets:])
			// checks for errors
			if err != nil {
				return 0, err
			}
			// if no bytes are read
			if n == 0 {
				if done {
					// If headers are complete, decide whether a body is expected.
					// If there's no Content-Length header, treat as done (no body).
					if _, clErr := r.Headers.Get("content-length"); clErr != nil {
						r.state = requestStateDone
					} else {
						r.state = requestStateParsingBody
					}
					// adds two bytes for CRLF
					totalByets += 2
				}
				break
			}
			if done {
				// If there is no content length switch to immediatly done
				// if not continue parsing body
				if _, clErr := r.Headers.Get("content-length"); clErr != nil {
					r.state = requestStateDone
				} else {
					r.state = requestStateParsingBody
				}
			}
			// increment total bytes by bytes consumed
			totalByets += n

		}
		// returns bytes consumed
		return totalByets, nil
	}

	// Parse Body
	if r.state == requestStateParsingBody {
		// gets the value of the content length header
		contentLengthStr, err := r.Headers.Get("content-length")
		//error check
		if err != nil {
			r.state = requestStateDone
			return 0, nil
		}
		//converts to int
		contentLength, err := strconv.Atoi(contentLengthStr)
		//error check
		if err != nil {
			return 0, errors.New("Invalid content length value")
		}
		//checks if thier still is a body
		//set's to done state if no body
		remaining := contentLength - len(r.Body)
		if remaining <= 0 {
			r.state = requestStateDone
			return 0, nil
		}
		//get the number of bytes read
		consumed := len(data)
		if consumed > remaining {
			consumed = remaining
		}
		// adds bytes read to body
		r.Body = append(r.Body, data[:consumed]...)
		// checks if done
		if len(r.Body) == contentLength {
			r.state = requestStateDone
		}
		return consumed, nil
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

// pareses request
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
		//Reades request returns bytes read
		request, err := reader.Read(chunck[parseTo:])
		//incrment bytes reaD
		parseTo += request
		// parses bytes and returns bytes parsed
		bytes, errored := httpRequest.parse(chunck[:parseTo])
		// reallocate chunck to account for the bytes read and bytes parsed
		copy(chunck, chunck[bytes:parseTo])
		// deincrment parse to by bytes that were parsed
		parseTo -= bytes
		// reset bytes parsed varible
		bytes = 0
		//Checks errors
		if errored != nil {
			return nil, errored
		}
		// when done reading
		if err == io.EOF {
			// if done
			if parseTo == 0 {
				switch httpRequest.state {
				case requestStateInitialized:
					// empty request
					return nil, errors.New("Empty")
				case requestStateDone:
					// done parsing
					break
				case requestStateParsingBody:
					// invaild content length given
					hasCL, clErr := httpRequest.Headers.Get("content-length")
					if clErr == nil {
						length, err := strconv.Atoi(hasCL)
						if err != nil {
							return nil, errors.New("Invalid content length value")
						}
						if len(httpRequest.Body) != length {
							return nil, errors.New("Incomplete body")
						}
					}
				default:
					// still parsing headers means missing "CRLF"
					return nil, errors.New("Incomplete headers")
				}
			} else if httpRequest.state == requestStateDone {
				// done parsing
				break
			} else if httpRequest.state == requestStateParsingBody {
				// same things as above if it is still parsing body then we were given invalid content length
				hasCL, clErr := httpRequest.Headers.Get("content-length")
				if clErr == nil {
					// convert to int
					length, err := strconv.Atoi(hasCL)
					if err != nil {
						return nil, errors.New("Invalid content length value")
					}
					if len(httpRequest.Body) != length {
						return nil, errors.New("Incomplete body")
					}
				}
				break
			} else if httpRequest.state == requestStateParsingHeaders {
				// still parsing headers means it's missing a CRLF
				return nil, errors.New("Incomplete headers")
			}
			break
		}
		if httpRequest.state == requestStateDone {
			// done parsing
			break
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
	// checks if method is uppercase
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
	// returns the request line and bytes read
	return Request, len(lines[0]) + 2, nil
}
