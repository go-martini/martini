// Package martini is a powerful package for quickly writing modular web applications/services in Golang.
//
// For a full guide visit http://github.com/go-martini/martini
//
//  package main
//
//  import "github.com/go-martini/martini"
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
	"crypto/tls"
	"github.com/codegangsta/inject"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"sync"
	"time"
)

// Martini represents the top level web application. inject.Injector methods can be invoked to map services on a global level.
type Martini struct {
	inject.Injector
	handlers      []Handler
	action        Handler
	httpListener  net.Listener
	httpsListener net.Listener
	mutex         sync.RWMutex
	logger        *log.Logger
}

// New creates a bare bones Martini instance. Use this method if you want to have full control over the middleware that is used.
func New() *Martini {
	m := &Martini{Injector: inject.New(), action: func() {}, logger: log.New(os.Stdout, "[martini] ", 0)}
	m.Map(m.logger)
	m.Map(defaultReturnHandler())
	return m
}

// Handlers sets the entire middleware stack with the given Handlers. This will clear any current middleware handlers.
// Will panic if any of the handlers is not a callable function
func (m *Martini) Handlers(handlers ...Handler) {
	m.handlers = make([]Handler, 0)
	for _, handler := range handlers {
		m.Use(handler)
	}
}

// Action sets the handler that will be called after all the middleware has been invoked. This is set to martini.Router in a martini.Classic().
func (m *Martini) Action(handler Handler) {
	validateHandler(handler)
	m.action = handler
}

// Logger sets the logger
func (m *Martini) Logger(logger *log.Logger) {
	m.logger = logger
	m.Map(m.logger)
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

// Run the http server on a given host and port.
func (m *Martini) RunOnAddr(addr string) {
	// TODO: Should probably be implemented using a new instance of http.Server in place of
	// calling http.ListenAndServer directly, so that it could be stored in the martini struct for later use.
	// This would also allow to improve testing when a custom host and port are passed.

	logger := m.Injector.Get(reflect.TypeOf(m.logger)).Interface().(*log.Logger)
	logger.Printf("listening on %s (%s)\n", addr, Env)
	//logger.Fatalln(http.ListenAndServe(addr, m))
	logger.Fatalln(m.listenAndServe(addr, m))
}

// Run the http server on a given host and port.
func (m *Martini) RunOnAddrTLS(addr, certFile, keyFile string) {
	// TODO: Should probably be implemented using a new instance of http.Server in place of
	// calling http.ListenAndServer directly, so that it could be stored in the martini struct for later use.
	// This would also allow to improve testing when a custom host and port are passed.

	logger := m.Injector.Get(reflect.TypeOf(m.logger)).Interface().(*log.Logger)
	logger.Printf("listening on %s (%s)\n", addr, Env)
	//logger.Fatalln(http.ListenAndServe(addr, m))
	logger.Fatalln(m.listenAndServeTLS(addr, certFile, keyFile, m))
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return
	}
	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

func (m *Martini) listenAndServe(addr string, handler http.Handler) error {
	server := &http.Server{Addr: addr, Handler: handler}

	if addr == "" {
		addr = ":http"
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	m.mutex.Lock()
	m.httpListener = ln
	m.mutex.Unlock()
	return server.Serve(tcpKeepAliveListener{ln.(*net.TCPListener)})
}

// listenAndServeTLS always returns a non-nil error.
func (m *Martini) listenAndServeTLS(addr, certFile, keyFile string, handler http.Handler) error {
	server := &http.Server{Addr: addr, Handler: handler}

	tlscfg := &tls.Config{}

	if tlscfg.NextProtos == nil {
		tlscfg.NextProtos = []string{"http/1.1"}
	}

	tlscfg.Certificates = make([]tls.Certificate, 1)
	var err error
	tlscfg.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return err
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	m.mutex.Lock()
	m.httpsListener = ln
	m.mutex.Unlock()

	tlsListener := tls.NewListener(tcpKeepAliveListener{ln.(*net.TCPListener)}, tlscfg)

	return server.Serve(tlsListener)

	//return server.ListenAndServeTLS(certFile, keyFile)
}

// Run the http server. Listening on os.GetEnv("PORT") or 3000 by default.
func (m *Martini) Run() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}

	host := os.Getenv("HOST")

	m.RunOnAddr(host + ":" + port)
}

// Run the http server. Listening on os.GetEnv("PORT") or 3000 by default.
func (m *Martini) RunTLS(certFile, keyFile string) {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "4000"
	}

	host := os.Getenv("HOST")

	m.RunOnAddrTLS(host+":"+port, certFile, keyFile)
}

func (m *Martini) Stop() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.httpListener != nil {
		m.httpListener.Close()
	}
	if m.httpsListener != nil {
		m.httpsListener.Close()
	}
}

func (m *Martini) createContext(res http.ResponseWriter, req *http.Request) *context {
	c := &context{inject.New(), m.handlers, m.action, NewResponseWriter(res), 0}
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

// Classic creates a classic Martini with some basic default middleware - martini.Logger, martini.Recovery and martini.Static.
// Classic also maps martini.Routes as a service.
func Classic() *ClassicMartini {
	r := NewRouter()
	m := New()
	m.Use(Logger())
	m.Use(Recovery())
	m.Use(Static("public"))
	m.MapTo(r, (*Routes)(nil))
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
	// Written returns whether or not the response for this context has been written.
	Written() bool
}

type context struct {
	inject.Injector
	handlers []Handler
	action   Handler
	rw       ResponseWriter
	index    int
}

func (c *context) handler() Handler {
	if c.index < len(c.handlers) {
		return c.handlers[c.index]
	}
	if c.index == len(c.handlers) {
		return c.action
	}
	panic("invalid index for context handler")
}

func (c *context) Next() {
	c.index += 1
	c.run()
}

func (c *context) Written() bool {
	return c.rw.Written()
}

func (c *context) run() {
	for c.index <= len(c.handlers) {
		_, err := c.Invoke(c.handler())
		if err != nil {
			panic(err)
		}
		c.index += 1

		if c.Written() {
			return
		}
	}
}
