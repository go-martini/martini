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
	m.MapTo(log.New(buff, "[martini] ", 0), (*Logger)(nil))
	m.Use(LoggerMiddleware())
	m.Use(func(res http.ResponseWriter) {
		res.WriteHeader(http.StatusNotFound)
	})

	req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusNotFound)
	refute(t, len(buff.String()), 0)
}
