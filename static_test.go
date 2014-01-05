package martini

import (
	"bytes"
	"github.com/codegangsta/inject"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Static(t *testing.T) {
	response := httptest.NewRecorder()

	m := New()

	m.Use(Static("."))

	req, err := http.NewRequest("GET", "http://localhost:3000/martini.go", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)
}

func Test_Static_As_Post(t *testing.T) {
	response := httptest.NewRecorder()

	m := New()
	r := NewRouter()

	m.Use(Static("."))
	m.Action(r.Handle)

	req, err := http.NewRequest("POST", "http://localhost:3000/martini.go", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusNotFound)
}

func Test_Static_BadDir(t *testing.T) {
	response := httptest.NewRecorder()

	m := Classic()

	req, err := http.NewRequest("GET", "http://localhost:3000/martini.go", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	refute(t, response.Code, http.StatusOK)
}

func Test_Static_Logging(t *testing.T) {
	response := httptest.NewRecorder()

	var buffer bytes.Buffer
	m := &Martini{inject.New(), []Handler{}, func() {}, log.New(&buffer, "[martini] ", 0)}
	m.Map(m.logger)
	m.Map(defaultReturnHandler())

	m.Use(Static("."))

	req, err := http.NewRequest("GET", "http://localhost:3000/martini.go", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)

	expect(t, buffer.String(), "[martini] [Static] Serving /martini.go\n") // Check logging output
}

func Test_Static_NoLogging(t *testing.T) {
	response := httptest.NewRecorder()

	var buffer bytes.Buffer
	m := &Martini{inject.New(), []Handler{}, func() {}, log.New(&buffer, "[martini] ", 0)}
	m.Map(m.logger)
	m.Map(defaultReturnHandler())

	m.Use(Static(".", true)) // Skip logging

	req, err := http.NewRequest("GET", "http://localhost:3000/martini.go", nil)
	if err != nil {
		t.Error(err)
	}

	m.ServeHTTP(response, req)
	expect(t, response.Code, http.StatusOK)

	expect(t, buffer.String(), "") // Check logging output
}
