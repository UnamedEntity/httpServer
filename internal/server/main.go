package server

import (
	"fmt"
	"net"
)

type Server struct {
	state    ServerState
	listener net.Listener
}

type ServerState int

const (
	initalize ServerState = iota
	parse
	done
)

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
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: 12\r\n\r\nHello World!")
	conn.Close()
}

func Serve(port int) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		return nil, err
	}
	server := &Server{
		initalize,
		listener,
	}
	go server.Listen()
	return server, nil

}
