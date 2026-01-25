package response

import (
	"io"
)


type StatusCode int

const (
	StatusOK 		  	StatusCode = 200
	StatusBadRequest		StatusCode = 400
	StatusInternalServerError	StatusCode = 500
)

func GetStatusLine(statusCode StatusCode) []byte {
	reasonPhrase := ""
	switch statusCode {
	case StatusOK:
		reasonPhrase = "OK"
	case StatusBadRequest:
		reasonPhrase = "Bad Request"
	case StatusInternalServerError:
		reasonPhrase = "Internal Server Error"
	return []byte(fmt.Sprintf("HTTP/1.1 %d %s\r\n", statusCode, reasonPhrase))
}

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {

}


