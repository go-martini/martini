package martini

import (
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
	m := New()
	refute(t, m, nil)
}

func Test_Martini_ServeHTTP(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

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

	m.ServeHTTP(response, (*http.Request)(nil))

	expect(t, result, "foobarbatbazban")
	expect(t, response.Code, 400)
}

func Test_Martini_EarlyWrite(t *testing.T) {
	result := ""
	response := httptest.NewRecorder()

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

	m.ServeHTTP(response, (*http.Request)(nil))

	expect(t, result, "foobar")
	expect(t, response.Code, 200)
}

func Test_Martini_UrlFor(t *testing.T) {
	m := Classic()

	m.Get("/foo", func() {
		// Nothing
	}).Name("foo_route")

	m.Post("/bar/:id", func(params Params) {
		// Nothing
	}).Name("bar_route")

	m.Post("/bar/:id/:name", func(params Params) {
		// Nothing
	}).Name("bar_id_name_route")

	expect(t, m.UrlFor("foo_route", nil), "/foo")
	expect(t, m.UrlFor("bar_route", 5), "/bar/5")
	expect(t, m.UrlFor("bar_id_name_route", 5, "john"), "/bar/5/john")
	expect(t, m.UrlFor("non_existent_route", nil), "")
}
