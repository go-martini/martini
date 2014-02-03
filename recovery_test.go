package martini

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Recovery(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()

	setENV(Dev)
	m := New()
	// replace log for testing
	m.Map(log.New(buff, "[martini] ", 0))
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "unpredictable")
	})
	m.Use(Recovery())
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		panic("here is a panic!")
	})
	m.ServeHTTP(recorder, (*http.Request)(nil))
	expect(t, recorder.Code, http.StatusInternalServerError)
	expect(t, recorder.HeaderMap.Get("Content-Type"), "text/html")
	refute(t, recorder.Body.Len(), 0)
	refute(t, len(buff.String()), 0)
}
