package martini

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Test_New(t *testing.T) {
	Convey("When a new Martini is created", t, func() {
		m := New()

		Convey("It should not be nil", func() {
			So(m, ShouldNotBeNil)
		})
	})
}

func Test_Martini_ServeHTTP(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	Convey("Given a brand new Martini with middleware and an action", t, func() {
		m := New()

		m.Use(func(c Context) {
			result += "foo"
			c.Next()
			result += "ban"
		})
		m.Use(func(c Context) {
			result += "bar"
			c.Next()
			result += "baz"
		})
		m.Action(func(res http.ResponseWriter, req *http.Request) {
			result += "bat"
			res.WriteHeader(400)
		})

		Convey("When the Martini serves HTTP", func() {
			m.ServeHTTP(response, (*http.Request)(nil))

			Convey("Middleware should have executed properly", func() {
				So(result, ShouldEqual, "foobarbatbazban")
			})
			Convey("The response code should be as expected", func() {
				So(response.Code, ShouldEqual, 400)
			})
		})
	})
}

func Test_Martini_EarlyWrite(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

	Convey("Given a Martini with middleware that writes out", t, func() {
		m := New()

		m.Use(func(res http.ResponseWriter) {
			result += "foobar"
			res.Write([]byte("Hello world"))
		})
		m.Use(func() {
			result += "bat"
		})
		m.Action(func(res http.ResponseWriter) {
			result += "baz"
			res.WriteHeader(400)
		})

		Convey("When HTTP is served", func() {
			m.ServeHTTP(response, (*http.Request)(nil))

			Convey("The middleware should execute as expected", func() {
				So(result, ShouldEqual, "foobar")
			})
			Convey("The response code should be 200, not 400", func() {
				So(response.Code, ShouldEqual, 200)
			})
			Convey("The response body should not be empty", func() {
				So(response.Body.String(), ShouldNotBeBlank)
			})
		})
	})
}
