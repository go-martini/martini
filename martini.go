package martini

import (
	"github.com/codegangsta/inject"
	"net/http"
	"reflect"
)

type Martini struct {
	handlers []interface{}
}

func New() *Martini {
	return &Martini{}
}

func (m *Martini) Use(handler interface{}) {
	m.handlers = append(m.handlers, handler)
}

func (m *Martini) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := &context{inject.New(), m.handlers, 0}
	ctx.MapTo(ctx, (*Context)(nil))
	ctx.MapTo(res, (*http.ResponseWriter)(nil))
	ctx.Map(req)
	ctx.run()
}

type Context interface {
	inject.Injector
	Next()
}

type context struct {
	injector inject.Injector
	handlers []interface{}
	index    int
}

func (c *context) Invoke(f interface{}) error {
	return c.injector.Invoke(f)
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
