package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
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
	request_line, httpVersion, err := ParseRequestLine(request)
	//Checks errors
	if err != nil {
		return nil, err
	}
	//fromats post request
	poseRequest := RequestLine{
		httpVersion[1],
		request_line[1],
		request_line[0],
	}
	httpRequest := &Request{
		poseRequest,
	}
	//return request
	return httpRequest, nil
}
func ParseRequestLine(f []byte) ([]string, []string, error) {
	//splits string into lines
	requestLines := strings.Split(string(f), "\r\n")
	//gets the request line
	requestLines = requestLines[0:1]
	// Splits the request line into the three parts
	request_lines := strings.Split(requestLines[0], " ")
	lowerCase := "abcdefghijklmnopqrstuvwxyz"
	// Checks to see if the request has the right input
	if len(request_lines) != 3 {
		return []string{}, []string{}, errors.New("Too many spaces")
	}
	//Checks for valid method
	for _, letter := range request_lines[0] {
		if strings.ContainsAny(string(letter), lowerCase) == true {
			return []string{}, []string{}, errors.New("Contains lower case")
		}
	}
	//Checks for right Http version
	if request_lines[2] != "HTTP/1.1" {
		return []string{}, []string{}, errors.New("Wrong Version of Http")
	}
	// Gets http version
	httpVersion := strings.Split(request_lines[2], "/")
	return request_lines, httpVersion, nil
}
