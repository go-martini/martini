package martini

import (
	"io"
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

	req5, err := http.NewRequest("GET", "http://localhost:3000/fez/this/should/match", nil)
	if err != nil {
		t.Error(err)
	}
	context5 := New().createContext(recorder, req5)

	req6, err := http.NewRequest("PUT", "http://localhost:3000/pop/blah/blah/blah/bap/foo/", nil)
	if err != nil {
		t.Error(err)
	}
	context6 := New().createContext(recorder, req6)

	req7, err := http.NewRequest("DELETE", "http://localhost:3000/wap//pow", nil)
	if err != nil {
		t.Error(err)
	}
	context7 := New().createContext(recorder, req7)

	req8, err := http.NewRequest("HEAD", "http://localhost:3000/wap//pow", nil)
	if err != nil {
		t.Error(err)
	}
	context8 := New().createContext(recorder, req8)

	req9, err := http.NewRequest("OPTIONS", "http://localhost:3000/opts", nil)
	if err != nil {
		t.Error(err)
	}
	context9 := New().createContext(recorder, req9)

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
	router.Get("/fez/**", func(params Params) {
		expect(t, params["_1"], "this/should/match")
		result += "fez"
	})
	router.Put("/pop/**/bap/:id/**", func(params Params) {
		expect(t, params["id"], "foo")
		expect(t, params["_1"], "blah/blah/blah")
		expect(t, params["_2"], "")
		result += "popbap"
	})
	router.Delete("/wap/**/pow", func(params Params) {
		expect(t, params["_1"], "")
		result += "wappow"
	})
	router.Options("/opts", func() {
		result += "opts"
	})
	router.Head("/wap/**/pow", func(params Params) {
		expect(t, params["_1"], "")
		result += "wappow"
	})

	router.Handle(recorder, req, context)
	router.Handle(recorder, req2, context2)
	router.Handle(recorder, req3, context3)
	router.Handle(recorder, req4, context4)
	router.Handle(recorder, req5, context5)
	router.Handle(recorder, req6, context6)
	router.Handle(recorder, req7, context7)
	router.Handle(recorder, req8, context8)
	router.Handle(recorder, req9, context9)
	expect(t, result, "foobarbatbarfoofezpopbapwappowwappowopts")
	expect(t, recorder.Code, http.StatusNotFound)
	expect(t, recorder.Body.String(), "404 page not found\n")
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
	handlers := []Handler{
		func() {
			result += "foo"
		},
		func(c Context) {
			result += "bar"
			c.Next()
			result += "bing"
		},
		func(resp http.ResponseWriter) {
			result += "bat"
			io.WriteString(resp, "Hello world")
		},
		func() {
			result += "baz"
		},
	}

	router.Get("/foo", handlers...)
	router.Handle(recorder, req, context)
	expect(t, result, "foobarbatbing")
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

func Test_NotFoundAsHandler(t *testing.T) {
	router := NewRouter()
	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "http://localhost:3000/foo", nil)
	context := New().createContext(recorder, req)

	router.Get("/bar", func() {})

	router.NotFound(func(resp http.ResponseWriter) {
		resp.WriteHeader(404)
		io.WriteString(resp, "not found")
	})

	router.Handle(recorder, req, context)
	expect(t, recorder.Code, http.StatusNotFound)
	expect(t, recorder.Body.String(), "not found")
}

func Test_NotFoundStacking(t *testing.T) {
	router := NewRouter()
	recorder := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "http://localhost:3000/foo", nil)
	if err != nil {
		t.Error(err)
	}
	context := New().createContext(recorder, req)

	result := ""
	handlers := []Handler{
		func() {
			result += "foo"
		},
		func(c Context) {
			result += "bar"
			c.Next()
			result += "bing"
		},
		func(resp http.ResponseWriter) {
			result += "bat"
			io.WriteString(resp, "Not Found")
		},
		func() {
			result += "baz"
		},
	}

	router.NotFound(handlers...)
	router.Handle(recorder, req, context)
	expect(t, result, "foobarbatbing")
	expect(t, recorder.Body.String(), "Not Found")
}

func Test_Any(t *testing.T) {
	router := NewRouter()
	router.Any("/foo", func(res http.ResponseWriter) {
		http.Error(res, "Nope", http.StatusNotFound)
	})

	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://localhost:3000/foo", nil)
	context := New().createContext(recorder, req)
	router.Handle(recorder, req, context)

	expect(t, recorder.Code, http.StatusNotFound)
	expect(t, recorder.Body.String(), "Nope\n")

	recorder = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "http://localhost:3000/foo", nil)
	context = New().createContext(recorder, req)
	router.Handle(recorder, req, context)

	expect(t, recorder.Code, http.StatusNotFound)
	expect(t, recorder.Body.String(), "Nope\n")
}

func Test_URLFor(t *testing.T) {
	router := NewRouter()
	var barIDNameRoute, fooRoute, barRoute Route

	fooRoute = router.Get("/foo", func() {
		// Nothing
	})

	barRoute = router.Post("/bar/:id", func(params Params) {
		// Nothing
	})

	barIDNameRoute = router.Get("/bar/:id/:name", func(params Params, routes Routes) {
		expect(t, routes.URLFor(fooRoute, nil), "/foo")
		expect(t, routes.URLFor(barRoute, 5), "/bar/5")
		expect(t, routes.URLFor(barIDNameRoute, 5, "john"), "/bar/5/john")
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
