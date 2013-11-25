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

	req4, err := http.NewRequest("PATCH", "http://localhost:3000/bar/foo", nil)
	if err != nil {
		t.Error(err)
	}
	context4 := New().createContext(recorder, req4)

	result := ""
	router.Get("/foo", func(req *http.Request) {
		result += "foo"
	})
	router.Patch("/bar/:id", func(params Params) {
		expect(t, params["id"], "foo")
		result += "barfoo"
	})
	router.Post("/bar/:id", func(params Params) {
		expect(t, params["id"], "bat")
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
	router.Handle(recorder, req4, context4)
	expect(t, result, "foobarbatbarfoo")
	expect(t, recorder.Code, http.StatusNotFound)
	expect(t, recorder.Body.String(), "404 page not found\n")
}

func Test_RouterHandlerStatusCode(t *testing.T) {
	router := NewRouter()
	router.Get("/foo", func() string {
		return "foo"
	})
	router.Get("/bar", func() (int, string) {
		return http.StatusForbidden, "bar"
	})
	router.Get("/baz", func() (string, string) {
		return "baz", "BAZ!"
	})

	// code should be 200 if none is returned from the handler
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3000/foo", nil)
	if err != nil {
		t.Error(err)
	}
	context := New().createContext(recorder, req)
	router.Handle(recorder, req, context)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, recorder.Body.String(), "foo")

	// if a status code is returned, it should be used
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "http://localhost:3000/bar", nil)
	if err != nil {
		t.Error(err)
	}
	context = New().createContext(recorder, req)
	router.Handle(recorder, req, context)
	expect(t, recorder.Code, http.StatusForbidden)
	expect(t, recorder.Body.String(), "bar")

	// shouldn't use the first returned value as a status code if not an integer
	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "http://localhost:3000/baz", nil)
	if err != nil {
		t.Error(err)
	}
	context = New().createContext(recorder, req)
	router.Handle(recorder, req, context)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, recorder.Body.String(), "baz")
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

func Test_RouteMatching(t *testing.T) {
	route := newRoute("GET", "/foo/:bar/bat/:baz", nil)
	for _, tt := range routeTests {
		ok, params := route.Match(tt.method, tt.path)
		if ok != tt.ok || params["bar"] != tt.params["bar"] || params["baz"] != tt.params["baz"] {
			t.Errorf("expected: (%v, %v) got: (%v, %v)", tt.ok, tt.params, ok, params)
		}
	}
}

func Test_NotFound(t *testing.T) {
	router := NewRouter()
	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "http://localhost:3000/foo", nil)
	context := New().createContext(recorder, req)

	router.NotFound(func(res http.ResponseWriter) {
		http.Error(res, "Nope", http.StatusNotFound)
	})

	router.Handle(recorder, req, context)
	expect(t, recorder.Code, http.StatusNotFound)
	expect(t, recorder.Body.String(), "Nope\n")
}

func Test_UrlFor(t *testing.T) {
	router := NewRouter()
    var barIdNameRoute, fooRoute, barRoute Route
    
	fooRoute = router.Get("/foo", func() {
		// Nothing
	})

	barRoute = router.Post("/bar/:id", func(params Params) {
		// Nothing
	})

	barIdNameRoute = router.Get("/bar/:id/:name", func(params Params, routes Routes) {
        expect(t, routes.UrlFor(fooRoute, nil), "/foo")
        expect(t, routes.UrlFor(barRoute, 5), "/bar/5")
        expect(t, routes.UrlFor(barIdNameRoute, 5, "john"), "/bar/5/john")
	})

	// code should be 200 if none is returned from the handler
	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://localhost:3000/bar/foo/bar", nil)
	if err != nil {
		t.Error(err)
	}
	context := New().createContext(recorder, req)
	router.Handle(recorder, req, context)
}
