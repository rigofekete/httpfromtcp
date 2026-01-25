package server

import (
	"net"
	"sync/atomic"
	"log"
	"fmt"

	"github.com/rigofekete/httpfromtcp/internal/response"
)


type Server struct {
	// closed indicates whether the resource has been closed.
	// Uses atomic.Bool to allow safe concurrent access from multiple goroutines.
	closed 	atomic.Bool
	listener	net.Listener
}


func Serve(port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{ 
		listener: l,
	}
	go s.listen()
	return s, nil
}

func (s *Server) listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.closed.Load() {
				return 
			}
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go s.handle(conn)
	}
}


// TODO continue handle refactor
func (s *Server) handle(conn net.Conn) error {	
	defer conn.Close()
	
	err := response.WriteStatusLine(conn, response.StatusOK)
	ir err != nil {
		return err
	}
	response.GetDefaultHeaders()


	// _, err := conn.Write([]byte("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\n\r\nHello World!\n"))
	if err != nil {
		return err
	}
	return nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

