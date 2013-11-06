package martini

import (
	"github.com/codegangsta/inject"
	"log"
	"net/http"
	"os"
	"reflect"
)

type Martini struct {
	inject.Injector
	handlers []Handler
}

func New() *Martini {
	m := &Martini{inject.New(), []Handler{}}
	m.Map(log.New(os.Stdout, "[martini] ", 0))
	return m
}

func (m *Martini) Use(handler Handler) {
	m.handlers = append(m.handlers, handler)
}

func (m *Martini) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := &context{inject.New(), m.handlers, 0}
	ctx.SetParent(m)
	ctx.MapTo(ctx, (*Context)(nil))
	ctx.MapTo(res, (*http.ResponseWriter)(nil))
	ctx.Map(req)
	ctx.run()
}

type Handler interface{}

type Context interface {
	inject.Injector
	Next()
}

type context struct {
	injector inject.Injector
	handlers []Handler
	index    int
}

func (c *context) Invoke(f interface{}) error {
	return c.injector.Invoke(f)
}

func (c *context) Apply(val interface{}) error {
	return c.injector.Apply(val)
}

func (c *context) Map(val interface{}) {
	c.injector.Map(val)
}

func (c *context) MapTo(val interface{}, ifacePtr interface{}) {
	c.injector.MapTo(val, ifacePtr)
}

func (c *context) Get(t reflect.Type) reflect.Value {
	return c.injector.Get(t)
}

func (c *context) SetParent(p inject.Injector) {
	c.injector.SetParent(p)
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
