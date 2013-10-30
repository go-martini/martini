package martini

import (
	"github.com/codegangsta/inject"
	"reflect"
)

type Context interface {
	inject.Injector
}

type context struct {
	injector inject.Injector
	index    int
}

func NewContext() Context {
	return &context{inject.New(), 0}
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

func (c *context) run(handlers []interface{}) {
}
