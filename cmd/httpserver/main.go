package main

import (
	"log"
	"os/signal"
	"syscall"
	"os"
	"io"

	"github.com/rigofekete/httpfromtcp/internal/server"
	"github.com/rigofekete/httpfromtcp/internal/request"
	"github.com/rigofekete/httpfromtcp/internal/response"
)

const port = 42069

func main() {
	server, err := server.Serve(handler, port)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}


func handler(w io.Writer, req *request.Request) *server.HandlerError {
	handlerResponse := &server.HandlerError{}
	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		handlerResponse.StatusCode = response.StatusBadRequest
		handlerResponse.Message = "Your problem is not my problem\n"
		return handlerResponse
	case "/myproblem":
		handlerResponse.StatusCode = response.StatusInternalServerError
		handlerResponse.Message = "Woopsie, my bad\n"
		return handlerResponse
	default:
		w.Write([]byte("All good, frfr\n"))
		return nil
	}
}

