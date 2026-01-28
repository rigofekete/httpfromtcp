package response

import (
	"io"
	"fmt"

	"github.com/rigofekete/httpfromtcp/internal/headers"
)

type writerState int

const (
	writerStateStatusLine writerState = iota
	writerStateHeaders
	writerStateBody
	writerStateTrailers
)

type Writer struct {
	writerState writerState
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writerState: writerStateStatusLine,
		writer: w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != writerStateStatusLine {
		return fmt.Errorf("cannot write status line in state %d", w.writerState)
	}
	defer func() { w.writerState = writerStateHeaders }()
	_, err := w.writer.Write(GetStatusLine(statusCode))
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != writerStateHeaders {
		return fmt.Errorf("cannot write headers in state %d", w.writerState)
	}
	defer func() { w.writerState = writerStateBody}()

	for k, v := range headers {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}
	totalBytesWritten := 0
	bytesWritten, err := w.writer.Write([]byte(fmt.Sprintf("%x\r\n", len(p))))
	// This is an alternate way to write the hexadecimal value to the writer
	// chunkSize := len(p)
	// bytesWritten, err := fmt.Fprint(w.writer, "%x\r\n", chunksize)
	if err != nil {
		return totalBytesWritten, err
	}
	totalBytesWritten += bytesWritten

	bytesWritten, err = w.writer.Write(p)
	if err != nil {
		return totalBytesWritten, err
	}
	totalBytesWritten += bytesWritten

	bytesWritten, err = w.writer.Write([]byte("\r\n"))
	if err != nil {
		return totalBytesWritten, err
	}
	totalBytesWritten += bytesWritten
	return totalBytesWritten, nil
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != writerStateBody {
		return 0, fmt.Errorf("cannot write body in state %d", w.writerState)
	}
	bytesWritten, err := w.writer.Write([]byte("0\r\n"))
	if err != nil {
		return bytesWritten, err
	}
	w.writerState = writerStateTrailers
	return bytesWritten, nil
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.writerState != writerStateTrailers {
		return fmt.Errorf("cannot write trailers in state %d", w.writerState)
	}
	defer func() { w.writerState = writerStateBody }()

	for k, v := range h {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}
