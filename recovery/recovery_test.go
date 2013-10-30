package recovery

import (
	"github.com/codegangsta/martini"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Recovery(t *testing.T) {
	m := martini.New()
	m.Use(New())
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		panic("here is a panic!")
	})
	m.ServeHTTP(httptest.NewRecorder(), (*http.Request)(nil))
}
