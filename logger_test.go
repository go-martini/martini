package martini

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Logger(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()

	m := New()
	// replace log for testing
	m.Map(log.New(buff, "[martini] ", 0))
	m.Use(Logger())
	m.Use(func(res http.ResponseWriter) {
		res.WriteHeader(404)
	})

	req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(recorder, req)
	expect(t, recorder.Code, 404)
	refute(t, len(buff.String()), 0)
}
