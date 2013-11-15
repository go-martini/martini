package auth

import (
	"encoding/base64"
	"github.com/codegangsta/martini"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_BasicAuth(t *testing.T) {
	m := martini.New()
	m.Use(Basic("foo", "bar"))
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("hello"))
	})

	r, _ := http.NewRequest("GET", "foo", nil)

	Convey("Martini should not serve unauthorized requests", t, func() {
		recorder := httptest.NewRecorder()
		m.ServeHTTP(recorder, r)

		So(recorder.Code, ShouldEqual, http.StatusUnauthorized)
		So(recorder.Body.String(), ShouldNotEqual, "hello")
	})

	Convey("Given basic authentication", t, func() {
		auth := "Basic " + base64.StdEncoding.EncodeToString([]byte("foo:bar"))
		r.Header.Set("Authorization", auth)

		Convey("Martini should authorize the authenticated request", func() {
			recorder := httptest.NewRecorder()
			m.ServeHTTP(recorder, r)

			So(recorder.Code, ShouldNotEqual, http.StatusUnauthorized)
			So(recorder.Body.String(), ShouldEqual, "hello")
		})
	})
}
