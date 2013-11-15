package martini

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Logger(t *testing.T) {
	Convey("Given a Martini with logger middleware on the stack", t, func() {
		buff := bytes.NewBufferString("")
		recorder := httptest.NewRecorder()
		m := New()

		// replace log for testing
		m.Map(log.New(buff, "[martini] ", 0))
		m.Use(Logger())
		m.Use(func(res http.ResponseWriter) {
			res.WriteHeader(404)
		})

		Convey("When a request is made to a route not found", func() {
			req, err := http.NewRequest("GET", "http://localhost:3000/foobar", nil)
			if err != nil {
				t.Error(err)
			}
			m.ServeHTTP(recorder, req)

			Convey("The proper error code should be returned", func() {
				So(recorder.Code, ShouldEqual, http.StatusNotFound)
			})

			Convey("The log buffer should not be empty", func() {
				So(buff.String(), ShouldNotBeBlank)
			})
		})
	})
}
