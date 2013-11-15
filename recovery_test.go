package martini

import (
	"bytes"
	. "github.com/smartystreets/goconvey/convey"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Recovery(t *testing.T) {
	buff := bytes.NewBufferString("")
	recorder := httptest.NewRecorder()

	Convey("Given a Martini with panicking middleware", t, func() {
		m := New()

		// replace log for testing
		m.Map(log.New(buff, "[martini] ", 0))
		m.Use(Recovery())
		m.Use(func(res http.ResponseWriter, req *http.Request) {
			panic("here is a panic!")
		})

		Convey("When HTTP is served", func() {
			m.ServeHTTP(recorder, (*http.Request)(nil))

			Convey("The response code should be Internal Server Error", func() {
				So(recorder.Code, ShouldEqual, http.StatusInternalServerError)
			})
			Convey("The log buffer should not be empty", func() {
				So(buff.String(), ShouldNotBeBlank)
			})
		})
	})
}
