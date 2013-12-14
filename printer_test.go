package martini

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_PrinterString(t *testing.T) {
	recorder := httptest.NewRecorder()

	m := New()
	m.Use(Printer())
	m.Use(func() string {
		return "foo"
	})
	m.ServeHTTP(recorder, (*http.Request)(nil))
	expect(t, recorder.Body.String(), "foo")
}

func Test_PrinterStatusCode(t *testing.T) {
	recorder := httptest.NewRecorder()

	m := New()
	m.Use(Printer())
	m.Use(func() int {
		return http.StatusForbidden
	})
	m.ServeHTTP(recorder, (*http.Request)(nil))
	expect(t, recorder.Code, http.StatusForbidden)
}

func Test_PrinterStringAndStatusCode(t *testing.T) {
	recorder := httptest.NewRecorder()

	m := New()
	m.Use(Printer())
	m.Use(func() (int, string) {
		return http.StatusForbidden, "foo"
	})
	m.ServeHTTP(recorder, (*http.Request)(nil))
	expect(t, recorder.Code, http.StatusForbidden)
	expect(t, recorder.Body.String(), "foo")
}

func Test_PrinterWithRouter(t *testing.T) {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://localhost:3000/foo", nil)

	r := NewRouter()
	m := New()
	m.Use(Printer())
	m.Action(r.Handle)

	r.Get("/foo", func() (int, string) {
		return http.StatusForbidden, "foo"
	})
	m.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusForbidden)
	expect(t, recorder.Body.String(), "foo")
}

func Test_PrinterBailEarly(t *testing.T) {
	recorder := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "http://localhost:3000/foo", nil)

	r := NewRouter()
	m := New()
	m.Use(Printer())
	m.Use(func() int {
		return http.StatusForbidden
	})
	m.Action(r.Handle)

	r.Get("/foo", func() int {
		return http.StatusOK
	})
	m.ServeHTTP(recorder, req)
	expect(t, recorder.Code, http.StatusForbidden)
}
