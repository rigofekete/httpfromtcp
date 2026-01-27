package server

import (
	"net"
	"sync/atomic"
	"log"
	"fmt"

	"github.com/rigofekete/httpfromtcp/internal/response"
	"github.com/rigofekete/httpfromtcp/internal/request"
)


type Handler func(w *response.Writer, req *request.Request)

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
	w := response.NewWriter(conn)
	req, err := request.RequestFromReader(conn)
	if err != nil {
		w.WriteStatusLine(response.StatusBadRequest)
		body := []byte(fmt.Sprintf("Error parsing request: %v", err))
		w.WriteHeaders(response.GetDefaultHeaders(len(body)))
		w.WriteBody(body)
		return
	}
	s.handler(w, req)
}


// func (s *Server) handle(conn net.Conn) {
// 	defer conn.Close()
// 	req, err := request.RequestFromReader(conn)
// 	if err != nil {
// 		hErr := &HandlerError {
// 			StatusCode: response.StatusBadRequest,
// 			Message: err.Error(),
// 		}
// 		hErr.Write(conn)
// 		return
// 	}
//
// 	buf := bytes.NewBuffer([]byte{})
// 	s.handler(buf, req)
// 	b := buf.Bytes()
// 	response.WriteStatusLine(conn, response.StatusOK)
// 	headers := response.GetDefaultHeaders(len(b))
// 	response.WriteHeaders(conn, headers)
// 	conn.Write(b)
// 	return
// }

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

