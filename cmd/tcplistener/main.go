package main

import (
	"fmt"
	"log"
	"net"


	"github.com/rigofekete/httpfromtcp/internal/request"
)

const port = ":42069"

func main() {
	listenerTCP, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening TCP traffic: %s\n", err.Error())
	}
	defer listenerTCP.Close()

	fmt.Println("listening for TCP traffic on", port)

	for {
		conn, err := listenerTCP.Accept()
		if err != nil {
			log.Fatalf("error listnening TCP traffic: %s\n", err.Error())
		}

		fmt.Println("connection has been accepted from", conn.RemoteAddr())

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error getting request from reader: %w", err)
		}
		fmt.Println("Request line:")
		fmt.Printf("- Method: %s\n", req.RequestLine.Method)
		fmt.Printf("- Target: %s\n", req.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", req.RequestLine.HttpVersion)
	}
}

