package martini

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Routing(t *testing.T) {
	router := NewRouter()
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost:3000/foo", nil)
	if err != nil {
		t.Error(err)
	}
	context := New().createContext(recorder, req)

	req2, err := http.NewRequest("POST", "http://localhost:3000/bar/bat", nil)
	if err != nil {
		t.Error(err)
	}
	context2 := New().createContext(recorder, req2)

	req3, err := http.NewRequest("DELETE", "http://localhost:3000/baz", nil)
	if err != nil {
		t.Error(err)
	}
	context3 := New().createContext(recorder, req3)

	result := ""
	router.Get("/foo", func(req *http.Request) {
		result += "foo"
	})
	router.Post("/bar/bat", func() {
		result += "barbat"
	})
	router.Put("/fizzbuzz", func() {
		result += "fizzbuzz"
	})
	router.Delete("/bazzer", func(c Context) {
		result += "baz"
	})

	router.Handle(recorder, req, context)
	router.Handle(recorder, req2, context2)
	router.Handle(recorder, req3, context3)
	expect(t, result, "foobarbat")
	expect(t, recorder.Code, http.StatusNotFound)
}

func Test_RouterHandlerStacking(t *testing.T) {
	router := NewRouter()
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost:3000/foo", nil)
	if err != nil {
		t.Error(err)
	}
	context := New().createContext(recorder, req)

	result := ""

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

	router.Get("/foo", f1, f2, f3, f4)

	router.Handle(recorder, req, context)
	expect(t, result, "foobarbat")
	expect(t, recorder.Body.String(), "Hello world")
}
