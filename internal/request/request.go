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
	request, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	request_line, httpVersion, err := ParseRequestLine(request)
	if err != nil {
		return nil, err
	}
	poseRequest := RequestLine{
		httpVersion[1],
		request_line[1],
		request_line[0],
	}
	httpRequest := &Request{
		poseRequest,
	}
	return httpRequest, nil
}
func ParseRequestLine(f []byte) ([]string, []string, error) {
	requestLines := strings.Split(string(f), "\r\n")
	requestLines = requestLines[0:1]
	request_lines := strings.Split(requestLines[0], " ")
	lowerCase := "abcdefghijklmnopqrstuvwxyz"
	if len(request_lines) != 3 {
		return []string{}, []string{}, errors.New("Too many spaces")
	}
	for _, letter := range request_lines[0] {
		if strings.ContainsAny(string(letter), lowerCase) == true {
			return []string{}, []string{}, errors.New("Contains lower case")
		}
	}
	if request_lines[2] != "HTTP/1.1" {
		return []string{}, []string{}, errors.New("Wrong Version of Http")
	}
	httpVersion := strings.Split(request_lines[2], "/")
	return request_lines, httpVersion, nil
}
