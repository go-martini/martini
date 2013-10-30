package martini

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_RecoveryHandler(t *testing.T) {
	recorder := httptest.NewRecorder()

	m := New()
	m.Use(RecoveryHandler())
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		panic("here is a panic!")
	})
	m.ServeHTTP(recorder, (*http.Request)(nil))
	// TODO verify that a log is written to
  expect(t, recorder.Code, 500)
}
