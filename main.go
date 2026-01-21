package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"io" 
	"strings"
)

const inputFilePath = "messages.txt"

func main() {
	file, err := os.Open(inputFilePath)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	// os.File implements io.ReadCloser interface, we can use it as arg
	chStrings := getLinesChannel(file)
	for line := range chStrings {
		fmt.Printf("read: %s\n", line)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	chStrings := make(chan string)

	go func() {
		defer f.Close()
		defer close(chStrings)

		currentLine := ""
		for {
			buffer := make([]byte, 8, 8)
			n, err := f.Read(buffer) 
			if err != nil {
				if currentLine != "" {
					chStrings <- currentLine
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
				chStrings <- currentLine + parts[i]
				currentLine = ""
			}
			currentLine += parts[len(parts)-1]
		}
	}()

	return chStrings
}

