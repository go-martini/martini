package martini

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_Routing(t *testing.T) {
	router := NewRouter()
	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "http://localhost:3000/foo", nil)
	context := New().createContext(recorder, req)

	req2, _ := http.NewRequest("POST", "http://localhost:3000/bar/bat", nil)
	context2 := New().createContext(recorder, req2)

	req3, _ := http.NewRequest("DELETE", "http://localhost:3000/baz", nil)
	context3 := New().createContext(recorder, req3)

	req4, _ := http.NewRequest("PATCH", "http://localhost:3000/bar/foo", nil)
	context4 := New().createContext(recorder, req4)

	req5, _ := http.NewRequest("GET", "http://localhost:3000/fez/this/should/match", nil)
	context5 := New().createContext(recorder, req5)

	req6, _ := http.NewRequest("PUT", "http://localhost:3000/pop/blah/blah/blah/bap/foo/", nil)
	context6 := New().createContext(recorder, req6)

	req7, _ := http.NewRequest("DELETE", "http://localhost:3000/wap//pow", nil)
	context7 := New().createContext(recorder, req7)

	req8, _ := http.NewRequest("HEAD", "http://localhost:3000/wap//pow", nil)
	context8 := New().createContext(recorder, req8)

	req9, _ := http.NewRequest("OPTIONS", "http://localhost:3000/opts", nil)
	context9 := New().createContext(recorder, req9)

	req10, _ := http.NewRequest("HEAD", "http://localhost:3000/foo", nil)
	context10 := New().createContext(recorder, req10)

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
	router.Handle(recorder, req10, context10)
	expect(t, result, "foobarbatbarfoofezpopbapwappowwappowoptsfoo")
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
	router.Get("/bytes", func() []byte {
		return []byte("Bytes!")
	})
	router.Get("/interface", func() interface{} {
		return "Interface!"
	})

	// code should be 200 if none is returned from the handler
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://localhost:3000/foo", nil)
	context := New().createContext(recorder, req)
	router.Handle(recorder, req, context)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, recorder.Body.String(), "foo")

	// if a status code is returned, it should be used
	recorder = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "http://localhost:3000/bar", nil)
	context = New().createContext(recorder, req)
	router.Handle(recorder, req, context)
	expect(t, recorder.Code, http.StatusForbidden)
	expect(t, recorder.Body.String(), "bar")

	// shouldn't use the first returned value as a status code if not an integer
	recorder = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "http://localhost:3000/baz", nil)
	context = New().createContext(recorder, req)
	router.Handle(recorder, req, context)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, recorder.Body.String(), "baz")

	// Should render bytes as a return value as well.
	recorder = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "http://localhost:3000/bytes", nil)
	context = New().createContext(recorder, req)
	router.Handle(recorder, req, context)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, recorder.Body.String(), "Bytes!")

	// Should render interface{} values.
	recorder = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "http://localhost:3000/interface", nil)
	context = New().createContext(recorder, req)
	router.Handle(recorder, req, context)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, recorder.Body.String(), "Interface!")
}

func Test_RouterHandlerStacking(t *testing.T) {
	router := NewRouter()
	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "http://localhost:3000/foo", nil)
	context := New().createContext(recorder, req)

	result := ""

	f1 := func() {
		result += "foo"
	}

	f2 := func(c Context) {
		result += "bar"
		c.Next()
		result += "bing"
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

func Test_MethodsFor(t *testing.T) {
	router := NewRouter()
	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "http://localhost:3000/foo", nil)
	context := New().createContext(recorder, req)
	router.Post("/foo/bar", func() {
	})

	router.Post("/fo", func() {
	})

	router.Get("/foo", func() {
	})

	router.Put("/foo", func() {
	})

	router.NotFound(func(routes Routes, w http.ResponseWriter, r *http.Request) {
		methods := routes.MethodsFor(r.URL.Path)
		if len(methods) != 0 {
			w.Header().Set("Allow", strings.Join(methods, ","))
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	router.Handle(recorder, req, context)
	expect(t, recorder.Code, http.StatusMethodNotAllowed)
	expect(t, recorder.Header().Get("Allow"), "GET,PUT")
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

	router.NotFound(func() string {
		return "not found"
	})

	router.Handle(recorder, req, context)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, recorder.Body.String(), "not found")

	recorder = httptest.NewRecorder()

	context = New().createContext(recorder, req)

	router.NotFound(func() (int, string) {
		return 404, "not found"
	})

	router.Handle(recorder, req, context)
	expect(t, recorder.Code, http.StatusNotFound)
	expect(t, recorder.Body.String(), "not found")

	recorder = httptest.NewRecorder()

	context = New().createContext(recorder, req)

	router.NotFound(func() (int, string) {
		return 200, ""
	})

	router.Handle(recorder, req, context)
	expect(t, recorder.Code, http.StatusOK)
	expect(t, recorder.Body.String(), "")
}

func Test_NotFoundStacking(t *testing.T) {
	router := NewRouter()
	recorder := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "http://localhost:3000/foo", nil)
	context := New().createContext(recorder, req)

	result := ""

	f1 := func() {
		result += "foo"
	}

	f2 := func(c Context) {
		result += "bar"
		c.Next()
		result += "bing"
	}

	f3 := func() string {
		result += "bat"
		return "Not Found"
	}

	f4 := func() {
		result += "baz"
	}

	router.NotFound(f1, f2, f3, f4)

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

	router.Get("/foo", func() {
		// Nothing
	}).Name("foo")

	router.Post("/bar/:id", func(params Params) {
		// Nothing
	}).Name("bar")

	router.Get("/bar/:id/:name", func(params Params, routes Routes) {
		expect(t, routes.URLFor("foo", nil), "/foo")
		expect(t, routes.URLFor("bar", 5), "/bar/5")
		expect(t, routes.URLFor("bar_id", 5, "john"), "/bar/5/john")
	}).Name("bar_id")

	// code should be 200 if none is returned from the handler
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://localhost:3000/bar/foo/bar", nil)
	context := New().createContext(recorder, req)
	router.Handle(recorder, req, context)
}
