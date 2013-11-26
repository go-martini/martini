package martini

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
)

// ResponseWriter is a wrapper around http.ResponseWriter that provides extra information about
// the response. It is recommended that middleware handlers use this construct to wrap a responsewriter
// if the functionality calls for it.
type ResponseWriter interface {
	http.ResponseWriter
	// Status returns the status code of the response or 0 if the response has not been written.
	Status() int
	// Written returns whether or not the ResponseWriter has been written.
	Written() bool
	// Size returns the size of the response body.
	Size() int
	// The ResponseWriter can no longer be written to.
	Close()
}

// NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter
func NewResponseWriter(rw http.ResponseWriter) ResponseWriter {
	return &responseWriter{rw, 0, 0, false}
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
	closed bool
}

func (rw *responseWriter) WriteHeader(s int) {
	if !rw.closed {
		rw.ResponseWriter.WriteHeader(s)
		rw.status = s
	}
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.closed {
		return 0, fmt.Errorf("ResponseWriter is now closed. Have you already written a response?")
	}
	if !rw.Written() {
		// The status will be StatusOK if WriteHeader has not been called yet
		rw.status = http.StatusOK
	}
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) Size() int {
	return rw.size
}

func (rw *responseWriter) Written() bool {
	return rw.status != 0
}

func (rw *responseWriter) Close() {
	rw.closed = true
}

func (rw *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hijacker, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("ResponseWriter doesn't support Hijacker interface")
	}
	return hijacker.Hijack()
}
