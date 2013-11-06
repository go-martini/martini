package martini

import (
	"errors"
	"github.com/codegangsta/inject"
	"log"
	"net/http"
	"os"
	"reflect"
)

type Martini interface {
	inject.Injector
	http.Handler
	Use(Handler) error
}

type martini struct {
	inject.Injector
	handlers []Handler
}

func New() Martini {
	m := &martini{inject.New(), []Handler{}}
	m.Map(log.New(os.Stdout, "[martini] ", 0))
	return m
}

func (m *martini) Use(handler Handler) error {
	if err := validateHandler(handler); err != nil {
		return err
	}

	m.handlers = append(m.handlers, handler)
	return nil
}

func (m *martini) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := &context{inject.New(), m.handlers, 0}
	ctx.SetParent(m)
	ctx.MapTo(ctx, (*Context)(nil))
	ctx.MapTo(res, (*http.ResponseWriter)(nil))
	ctx.Map(req)
	ctx.run()
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
