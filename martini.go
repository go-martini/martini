package martini

import (
	"github.com/codegangsta/inject"
	"log"
	"net/http"
	"os"
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
