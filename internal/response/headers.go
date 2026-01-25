package response

import (
	"fmt"
	"io"

	"github.com/rigofekete/httpfromtcp/internal/headers"
)

func GetDefaultHeaders(contentLen int) headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}


func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for k, v := range headers {
		fieldLine := fmt.Sprintf("%s: %s\r\n", k, v)
		_, err := w.Write([]byte(fieldLine))
		if err != nil {
			return err
		}
	}
	_, err := w.Write([]byte("\r\n"))
	return err
}



