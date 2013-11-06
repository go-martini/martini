package martini

import (
	"errors"
	"github.com/codegangsta/inject"
	"log"
	"net/http"
	"os"
	"reflect"
)

type Martini struct {
	inject.Injector
	handlers []Handler
	action   Handler
}

func New() *Martini {
	m := &Martini{inject.New(), []Handler{}, func() {}}
	m.Map(log.New(os.Stdout, "[martini] ", 0))
	return m
}

func (m *Martini) Use(handler Handler) error {
	if err := validateHandler(handler); err != nil {
		return err
	}

	m.handlers = append(m.handlers, handler)
	return nil
}

func (m *Martini) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	m.createContext(res, req).run()
}

func (m *Martini) Action(handler Handler) error {
	if err := validateHandler(handler); err != nil {
		return err
	}

	m.action = handler
	return nil
}

func (m *Martini) createContext(res http.ResponseWriter, req *http.Request) *context {
	c := &context{inject.New(), append(m.handlers, m.action), 0}
	c.SetParent(m)
	c.MapTo(c, (*Context)(nil))
	c.MapTo(res, (*http.ResponseWriter)(nil))
	c.Map(req)
	return c
}

type ClassicMartini struct {
	*Martini
	Router
}

func Classic() *ClassicMartini {
	r := NewRouter()
	m := New()
	m.Use(Logger)
	m.Use(Recovery)
	m.Action(r.Handle)
	return &ClassicMartini{m, r}
}

type Handler interface{}

func validateHandler(handler Handler) error {
	if reflect.TypeOf(handler).Kind() != reflect.Func {
		return errors.New("martini handler must be a callable func")
	}
	return nil
}

type Context interface {
	inject.Injector
	Next()
}

type context struct {
	inject.Injector
	handlers []Handler
	index    int
}

func (c *context) Next() {
	c.index += 1
	c.run()
}

func (c *context) run() {
	for c.index < len(c.handlers) {
		err := c.Invoke(c.handlers[c.index])
		if err != nil {
			panic(err)
		}
		c.index += 1
	}
}
