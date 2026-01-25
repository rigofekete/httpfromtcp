package response

import (
	"fmt"

	"github.com/rigofekete/httpfromtcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{
		"Content-Length": 	fmt.Sprintf("%d", contentLen),
		"Connection": 		"close",
		"Content-Type":		"text/plain",
	}
}


func WriteHeaders(w io.Writer, headers headers.Headers) error {
	// TODO user headers.Get
	for k, v := range headers {
		fieldLine := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.Write([]byte(fieldLine))
		if err != nil {
			return err
		}
	}
	return err
}



