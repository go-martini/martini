package martini

import (
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Routing(t *testing.T) {
	recorder := httptest.NewRecorder()
	result := ""

	Convey("Given a router with some routes and handlers", t, func() {
		router := NewRouter()
		router.Get("/foo", func(req *http.Request) {
			result += "foo"
		})
		router.Post("/bar/:id", func(params Params) {
			So(params["id"], ShouldEqual, "bat")
			result += "barbat"
		})
		router.Put("/fizzbuzz", func() {
			result += "fizzbuzz"
		})
		router.Delete("/bazzer", func(c Context) {
			result += "baz"
		})

		Convey("And a series of requests to route", func() {
			req, err := http.NewRequest("GET", "http://localhost:3000/foo", nil)
			So(err, ShouldBeNil)
			context := New().createContext(recorder, req)

			req2, err := http.NewRequest("POST", "http://localhost:3000/bar/bat", nil)
			So(err, ShouldBeNil)
			context2 := New().createContext(recorder, req2)

			req3, err := http.NewRequest("DELETE", "http://localhost:3000/baz", nil)
			So(err, ShouldBeNil)
			context3 := New().createContext(recorder, req3)

			Convey("When the requests are handled in order", func() {
				router.Handle(recorder, req, context)
				router.Handle(recorder, req2, context2)
				router.Handle(recorder, req3, context3)

				Convey("The request should be routed and handled correctly", func() {
					So(result, ShouldEqual, "foobarbat")
					So(recorder.Code, ShouldEqual, http.StatusNotFound)
					So(recorder.Body.String(), ShouldEqual, http.StatusText(http.StatusNotFound))
				})
			})
		})
	})
}

func Test_RouterHandlerStacking(t *testing.T) {
	recorder := httptest.NewRecorder()
	result := ""

	Convey("With a router and set of handlers", t, func() {
		router := NewRouter()

		req, err := http.NewRequest("GET", "http://localhost:3000/foo", nil)
		So(err, ShouldBeNil)
		context := New().createContext(recorder, req)

		f1 := func() {
			result += "foo"
		}
		f2 := func() {
			result += "bar"
		}
		f3 := func() string {
			result += "bat"
			return "Hello world"
		}
		f4 := func() {
			result += "baz"
		}

		Convey("Stacked route should be handled as expected", func() {
			router.Get("/foo", f1, f2, f3, f4)
			router.Handle(recorder, req, context)

			So(result, ShouldEqual, "foobarbat")
			So(recorder.Body.String(), ShouldEqual, "Hello world")
		})
	})
}

func Test_RouteMatching(t *testing.T) {
	Convey("With a route to match", t, func() {
		route := newRoute("GET", "/foo/:bar/bat/:baz", nil)

		Convey("Each route string should match the correct route", func() {
			for _, tt := range routeTests {
				ok, params := route.match(tt.method, tt.path)
				So(ok, ShouldEqual, tt.ok)
				So(params["bar"], ShouldEqual, tt.params["bar"])
				So(params["baz"], ShouldEqual, tt.params["baz"])
			}
		})
	})
}

var routeTests = []struct {
	// in
	method string
	path   string

	// out
	ok     bool
	params map[string]string
}{
	{"GET", "/foo/123/bat/321", true, map[string]string{"bar": "123", "baz": "321"}},
	{"POST", "/foo/123/bat/321", false, map[string]string{}},
	{"GET", "/foo/hello/bat/world", true, map[string]string{"bar": "hello", "baz": "world"}},
	{"GET", "foo/hello/bat/world", false, map[string]string{}},
	{"GET", "/foo/123/bat/321/", true, map[string]string{"bar": "123", "baz": "321"}},
	{"GET", "/foo/123/bat/321//", false, map[string]string{}},
	{"GET", "/foo/123//bat/321/", false, map[string]string{}},
}
