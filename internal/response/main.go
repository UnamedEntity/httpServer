package response

import (
	"fmt"
	"httpServer/internal/headers"
	"io"
	"strconv"
)

type StatusCode int

const (
	Code200 StatusCode = 200
	Code400 StatusCode = 400
	Code500 StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	var err error

	switch statusCode {
	case Code200:
		_, err = fmt.Fprint(w, "HTTP/1.1 200 OK\r\n")
	case Code400:
		_, err = fmt.Fprint(w, "HTTP/1.1 400 Bad Request\r\n")
	case Code500:
		_, err = fmt.Fprint(w, "HTTP/1.1 500 Internal Server Error\r\n")
	default:
		_, err = fmt.Fprint(w, "HTTP/1.1 500\r\n")
	}
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	headers := headers.NewHeaders()
	headers["content-length"] = strconv.Itoa(contentLen)
	headers["connection"] = "close"
	headers["content-type"] = "text/plain"
	return headers
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for i, x := range headers {
		_, err := fmt.Fprintf(w, "%s: %v\r\n", i, x)
		if err != nil {
			return err
		}
	}
	fmt.Fprint(w, "\r\n")
	return nil
}

type Writer struct {
	write io.Writer
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	err := WriteStatusLine(w.write, statusCode)
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	err := WriteHeaders(w.write, headers)
	if err != nil {
		return err
	}
	return nil
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	length, err := fmt.Fprintf(w.write, "%s", string(p))
	if err != nil {
		return 0, err
	}
	return length, nil
}
