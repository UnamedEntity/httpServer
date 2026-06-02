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
		// accepts a conncetion from the listner
		file, err := s.listener.Accept()
		if err != nil {
			return
		}
		// sends it to the handle function asyncronously
		go s.handle(file)
	}
}
func (s *Server) handle(conn net.Conn) {
	// closes connection when function is done
	defer conn.Close()
	// parses the request and returns request struct
	req, err := request.RequestFromReader(conn)
	if err != nil {
		return
	}
	//creates a buffer
	buf := &bytes.Buffer{}
	// create a writer that writes to buffer
	respW := response.NewWriter(buf)
	// let the handler write status, headers, and body into respW
	s.handler(respW, req)
	// write buffered response to the connection
	_, err = conn.Write(buf.Bytes())
	// checks for errors
	if err != nil {
		return
	}
}

// Serve function to give a response to a request
func Serve(port int, handler Handler) (*Server, error) {
	//create a listner
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	//check for errors
	if err != nil {
		return nil, err
	}
	// assigns
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
