package recovery

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Recovery(t *testing.T) {
	recorder := httptest.NewRecorder()

	m := martini.New()
	m.Use(New())
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		panic("here is a panic!")
	})
	m.ServeHTTP(recorder, (*http.Request)(nil))
	// TODO verify that a log is written to
	if recorder.Code != 500 {
		t.Error("Response did not return 500")
	}
}
