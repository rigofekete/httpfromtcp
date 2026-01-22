package main

import (
	"errors"
	"fmt"
	"log"
	"io" 
	"strings"
	"net"
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

		linesCh := getLinesChannel(conn)

		for line := range linesCh {
			fmt.Println(line)
		}
		fmt.Println("connection to ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesCh := make(chan string)

	go func() {
		defer f.Close()
		defer close(linesCh)

		currentLine := ""
		for {
			buffer := make([]byte, 8, 8)
			n, err := f.Read(buffer) 
			if err != nil {
				if currentLine != "" {
					linesCh <- currentLine
					currentLine = ""
				}
				if errors.Is(err, io.EOF) {
					return 
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}

			str := string(buffer[:n])
			parts := strings.Split(str, "\n")
			for i := 0; i < len(parts)-1; i++ {
				linesCh <- fmt.Sprintf("%s%s", currentLine, parts[i])
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()

	return linesCh
}

