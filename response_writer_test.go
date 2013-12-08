package martini

import (
	"bufio"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type closeNotifyingRecorder struct {
	*httptest.ResponseRecorder
	closed chan bool
}

func newCloseNotifyingRecorder() *closeNotifyingRecorder {
	return &closeNotifyingRecorder{
		httptest.NewRecorder(),
		make(chan bool, 1),
	}
}

func (c *closeNotifyingRecorder) close() {
	c.closed <- true
}

func (c *closeNotifyingRecorder) CloseNotify() <-chan bool {
	return c.closed
}

type hijackableResponse struct {
	Hijacked bool
}

func newHijackableResponse() *hijackableResponse {
	return &hijackableResponse{}
}

func (h *hijackableResponse) Header() http.Header           { return nil }
func (h *hijackableResponse) Write(buf []byte) (int, error) { return 0, nil }
func (h *hijackableResponse) WriteHeader(code int)          {}
func (h *hijackableResponse) Flush()                        {}
func (h *hijackableResponse) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h.Hijacked = true
	return nil, nil, nil
}

func Test_ResponseWriter_WritingString(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	rw.Write([]byte("Hello world"))

	expect(t, rec.Code, rw.Status())
	expect(t, rec.Body.String(), "Hello world")
	expect(t, rw.Status(), http.StatusOK)
	expect(t, rw.Size(), 11)
	expect(t, rw.Written(), true)
}

func Test_ResponseWriter_WritingStrings(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	rw.Write([]byte("Hello world"))
	rw.Write([]byte("foo bar bat baz"))

	expect(t, rec.Code, rw.Status())
	expect(t, rec.Body.String(), "Hello worldfoo bar bat baz")
	expect(t, rw.Status(), http.StatusOK)
	expect(t, rw.Size(), 26)
}

func Test_ResponseWriter_WritingHeader(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)

	rw.WriteHeader(http.StatusNotFound)

	expect(t, rec.Code, rw.Status())
	expect(t, rec.Body.String(), "")
	expect(t, rw.Status(), http.StatusNotFound)
	expect(t, rw.Size(), 0)
}

func Test_ResponseWriter_Before(t *testing.T) {
	rec := httptest.NewRecorder()
	rw := NewResponseWriter(rec)
	result := ""

	rw.Before(func(ResponseWriter) {
		result += "foo"
	})
	rw.Before(func(ResponseWriter) {
		result += "bar"
	})

	rw.WriteHeader(http.StatusNotFound)

	expect(t, rec.Code, rw.Status())
	expect(t, rec.Body.String(), "")
	expect(t, rw.Status(), http.StatusNotFound)
	expect(t, rw.Size(), 0)
	expect(t, result, "barfoo")
}

func Test_ResponseWriter_Hijack(t *testing.T) {
	hijackable := newHijackableResponse()
	rw := NewResponseWriter(hijackable)
	hijacker, ok := rw.(http.Hijacker)
	expect(t, ok, true)
	_, _, err := hijacker.Hijack()
	if err != nil {
		t.Error(err)
	}
	expect(t, hijackable.Hijacked, true)
}

func Test_ResponseWrite_Hijack_NotOK(t *testing.T) {
	hijackable := new(http.ResponseWriter)
	rw := NewResponseWriter(*hijackable)
	hijacker, ok := rw.(http.Hijacker)
	expect(t, ok, true)
	_, _, err := hijacker.Hijack()

	refute(t, err, nil)
}

func Test_ResponseWriter_CloseNotify(t *testing.T) {
	rec := newCloseNotifyingRecorder()
	rw := NewResponseWriter(rec)
	closed := false
	notifier := rw.(http.CloseNotifier).CloseNotify()
	rec.close()
	select {
	case <-notifier:
		closed = true
	case <-time.After(time.Second):
	}
	expect(t, closed, true)
}
