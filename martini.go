// Martini is a powerful package for quickly writing modular web applications/services in Golang.
//
// For a full guide visit http://github.com/codegangsta/martini
//
//  package main
//
//  import "github.com/codegangsta/martini"
//
//  func main() {
//    m := martini.Classic()
//
//    m.Get("/", func() string {
//      return "Hello world!"
//    })
//
//    m.Run()
//  }
package martini

import (
	"github.com/codegangsta/inject"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
)

// Martini represents the top level web application. inject.Injector methods can be invoked to map services on a global level.
type Martini struct {
	inject.Injector
	handlers []Handler
	action   Handler
	logger   *log.Logger
}

// New creates a bare bones Martini instance. Use this method if you want to have full control over the middleware that is used.
func New() *Martini {
	m := &Martini{inject.New(), []Handler{}, func() {}, log.New(os.Stdout, "[martini] ", 0)}
	m.Map(m.logger)
	return m
}

// Use adds a middleware Handler to the stack. Will panic if the handler is not a callable func. Middleware Handlers are invoked in the order that they are added.
func (m *Martini) Use(handler Handler) {
	validateHandler(handler)

	m.handlers = append(m.handlers, handler)
}

// ServeHTTP is the HTTP Entry point for a Martini instance. Useful if you want to control your own HTTP server.
func (m *Martini) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	m.createContext(res, req).run()
}

// Action sets the handler that will be called after all the middleware has been invoked. This is set to martini.Router in a martini.Classic().
func (m *Martini) Action(handler Handler) {
	validateHandler(handler)
	m.action = handler
}

type RunConfig struct {
	Addr string
	Port int
}

var defaultConfig = &RunConfig{Port: 3000, Addr: "0.0.0.0"}

// Run the http server. Listening on os.GetEnv("PORT") or 3000 by default.
// You can use a martini.RunConfig to use a customized run config.
func (m *Martini) Run(config ...*RunConfig) {
	// Martini config management
	var martiniConfig *RunConfig
	if len(config) == 0 {
		// No RunConfig given: we use the default config, and see if
		// Env variables "PORT" and "ADDR" has been set.
		martiniConfig = defaultConfig
		addr := os.Getenv("ADDR")
		if len(addr) > 0 {
			martiniConfig.Addr = addr
		}
		port, err := strconv.Atoi(os.Getenv("PORT"))
		if err != nil && port > 0 {
			martiniConfig.Port = port
		}
	} else {
		// A RunConfig has been provided: we'll use that one!
		martiniConfig = config[0]
	}
	port := strconv.Itoa(martiniConfig.Port)
	addr := martiniConfig.Addr

	m.logger.Println("listening on address '" + addr + "', port " + port)
	http.ListenAndServe(addr+":"+port, m)
}

func (m *Martini) createContext(res http.ResponseWriter, req *http.Request) *context {
	c := &context{inject.New(), append(m.handlers, m.action), &responseWriter{res, false}, 0}
	c.SetParent(m)
	c.MapTo(c, (*Context)(nil))
	c.MapTo(c.rw, (*http.ResponseWriter)(nil))
	c.Map(req)
	return c
}

// ClassicMartini represents a Martini with some reasonable defaults. Embeds the router functions for convenience.
type ClassicMartini struct {
	*Martini
	Router
}

// Classic creates a classic Martini with some basic default middleware - martini.Logger, martini.Recovery, and martini.Static.
func Classic() *ClassicMartini {
	r := NewRouter()
	m := New()
	m.Use(Logger())
	m.Use(Recovery())
	m.Use(Static("public"))
	m.Action(r.Handle)
	return &ClassicMartini{m, r}
}

// Handler can be any callable function. Martini attempts to inject services into the handler's argument list.
// Martini will panic if an argument could not be fullfilled via dependency injection.
type Handler interface{}

func validateHandler(handler Handler) {
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		panic("martini handler must be a callable func")
	}
}

// Context represents a request context. Services can be mapped on the request level from this interface.
type Context interface {
	inject.Injector
	// Next is an optional function that Middleware Handlers can call to yield the until after
	// the other Handlers have been executed. This works really well for any operations that must
	// happen after an http request
	Next()
	written() bool
}

type context struct {
	inject.Injector
	handlers []Handler
	rw       *responseWriter
	index    int
}

func (c *context) Next() {
	c.index += 1
	c.run()
}

func (c *context) written() bool {
	return c.rw.written
}

func (c *context) run() {
	for c.index < len(c.handlers) {
		_, err := c.Invoke(c.handlers[c.index])
		if err != nil {
			panic(err)
		}
		c.index += 1

		if c.rw.written {
			return
		}
	}
}

type responseWriter struct {
	w       http.ResponseWriter
	written bool
}

func (r *responseWriter) Header() http.Header {
	return r.w.Header()
}

func (r *responseWriter) Write(b []byte) (int, error) {
	r.written = true
	return r.w.Write(b)
}

func (r *responseWriter) WriteHeader(s int) {
	r.written = true
	r.w.WriteHeader(s)
}
