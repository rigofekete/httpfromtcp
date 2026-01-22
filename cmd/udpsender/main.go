package main

import (
	"fmt"
	"net"
	"bufio"
	"os"
)

const serverAddr = "localhost:42069"

func main() {
	udpAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving UDP serverAddr: %v\n", err)
		os.Exit(1)
	}

	udpConn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error dialing UDP: %v\n", err)
		os.Exit(1)
	}
	defer udpConn.Close()

	fmt.Printf("Sending to %s. Type your message and press Enter to send. Press Ctrl+X to exit.\n", serverAddr)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}
		_, err = udpConn.Write([]byte(line))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error sending message: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Message sent: %s", line)
	}
}
