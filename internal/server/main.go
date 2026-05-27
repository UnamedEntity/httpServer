package server

import (
	"bytes"
	"fmt"
	"httpServer/internal/request"
	"httpServer/internal/response"
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
	respW := response.NewWriter(buf)
	// let the handler write status, headers, and body into respW
	s.handler(respW, req)
	// write buffered response to the connection
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

type Handler func(w *response.Writer, req *request.Request)
