package server

import (
	"net"
	"sync/atomic"
	"bytes"
	"log"
	"fmt"

	"github.com/rigofekete/httpfromtcp/internal/response"
	"github.com/rigofekete/httpfromtcp/internal/request"
)


type HandlerError struct {
	StatusCode response.StatusCode
	Message string
}

func (h *HandlerError) Write(conn io.Writer) {
	response.WriteStatusLine(conn, h.StatusCode)
	messageBytes := []byte(h.Message)
	headers := response.GetDefaultHeaders(len(messageBytes))
	response.WriteHeaders(conn, headers)
	conn.Write(messageBytes)
}

type Handler func(w io.Writer, req *request.Request) *HandlerError

// Server is an HTTP 1.1 server
type Server struct {
	// closed indicates whether the resource has been closed.
	// Uses atomic.Bool to allow safe concurrent access from multiple goroutines.
	closed 	atomic.Bool
	listener	net.Listener
	handler	Handler
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

func (s *Server) handle(conn net.Conn) {
	defer conn.Close()
	req, err := request.RequestFromReader(conn)
	if err != nil {
		hErr := &HandlerError {
			StatusCode: response.StatusBadRequest,
			Message: err.Error(),
		}
		hErr.Write(conn)
		return
	}

	buf := bytes.NewBuffer([]byte{})
	hErr := s.handler(buf, req)
	if hErr != nil {
		hErr.Write(conn)
		return
	}
	b := buf.Bytes()
	response.WriteStatusLine(conn, response.StatusOK)
	headers := response.GetDefaultHeaders(len(b))
	response.WriteHeaders(conn, headers)
	conn.Write(b)
	return
}

func Serve(handler Handler, port int) (*Server, error) {
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		listener: l,
		handler: handler,
	}
	go s.listen()
	return s, nil
}

func (s *Server) Close() error {
	s.closed.Store(true)
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

