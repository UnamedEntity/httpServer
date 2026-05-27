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

// Writer wraps an io.Writer and exposes helpers to write a full HTTP response
// in the required order: StatusLine -> Headers -> Body.
type Writer struct {
	write io.Writer
	state int
}

const (
	writerStateInit = iota
	writerStateStatusWritten
	writerStateHeadersWritten
)

// NewWriter returns a new response.Writer wrapping the provided io.Writer.
func NewWriter(w io.Writer) *Writer {
	return &Writer{write: w, state: writerStateInit}
}

// WriteStatusLine writes the HTTP status line. Must be called first.
func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.state != writerStateInit {
		return fmt.Errorf("invalid write order: status line already written")
	}
	if err := WriteStatusLine(w.write, statusCode); err != nil {
		return err
	}
	w.state = writerStateStatusWritten
	return nil
}

// WriteHeaders writes headers. Must be called after WriteStatusLine.
func (w *Writer) WriteHeaders(h headers.Headers) error {
	if w.state != writerStateStatusWritten {
		return fmt.Errorf("invalid write order: headers must follow status line")
	}
	if err := WriteHeaders(w.write, h); err != nil {
		return err
	}
	w.state = writerStateHeadersWritten
	return nil
}

// WriteBody writes the response body. Must be called after WriteHeaders.
func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.state != writerStateHeadersWritten {
		return 0, fmt.Errorf("invalid write order: body must follow headers")
	}
	return w.write.Write(p)
}
