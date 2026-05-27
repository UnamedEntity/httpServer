package server

import (
	"bytes"
	"fmt"
	"httpServer/internal/headers"
	"httpServer/internal/request"
	"httpServer/internal/response"
	"io"
	"net"
)

type Server struct {
	state    string
	listener net.Listener
	handler  Handler
}

func (s *Server) Close() error {
	s.listener.Close()
	return nil
}

func (s *Server) Listen() {
	for {
		file, err := s.listener.Accept()
		if err != nil {
			return
		}
		go s.handle(file)
	}
}
func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		return
	}
	buf := &bytes.Buffer{}
	handlerError := s.handler(buf, req)
	if handlerError != nil {
		handlerError.Write(conn)
		return
	}
	err = response.WriteStatusLine(conn, response.Code200)
	if err != nil {
		return
	}
	headers := response.GetDefaultHeaders(buf.Len())
	err = response.WriteHeaders(conn, headers)
	if err != nil {
		return
	}
	_, err = conn.Write(buf.Bytes())
	if err != nil {
		return
	}
}

// Serve function to give a response to a request
func Serve(port int, handler Handler) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		"",
		listener,
		handler,
	}
	go server.Listen()
	return server, nil

}

// Handeler types

type HandlerError struct {
	state   response.StatusCode
	message string
	headers headers.Headers
}

func (h *HandlerError) Write(w io.Writer) error {
	err := response.WriteStatusLine(w, h.state)
	if err != nil {
		return err
	}
	if h.headers != nil {
		err = response.WriteHeaders(w, h.headers)
	} else {
		headers := response.GetDefaultHeaders(len(h.message))
		err = response.WriteHeaders(w, headers)
	}
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(w, "%s", h.message)
	if err != nil {
		return err
	}
	return nil
}

// NewHandlerError creates a HandlerError with the given status and message.
func NewHandlerError(code response.StatusCode, msg string) *HandlerError {
	return &HandlerError{state: code, message: msg}
}

// NewHandlerErrorWithHeaders creates a HandlerError with custom headers.
func NewHandlerErrorWithHeaders(code response.StatusCode, msg string, hdr headers.Headers) *HandlerError {
	return &HandlerError{state: code, message: msg, headers: hdr}
}

type Handler func(w io.Writer, req *request.Request) *HandlerError
