package martini

import (
	"github.com/codegangsta/inject"
	"net/http"
	"reflect"
)

type Context interface {
	inject.Injector
	Next()
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

type ReturnHandler func(Context, []reflect.Value)

// Returns the matching handler
func defaultReturnHandler() ReturnHandler {
	return func(ctx Context, vals []reflect.Value) {
		rv := ctx.Get(inject.InterfaceOf((*http.ResponseWriter)(nil)))
		res := rv.Interface().(http.ResponseWriter)
		var responseVal reflect.Value
		if len(vals) > 1 && vals[0].Kind() == reflect.Int {
			res.WriteHeader(int(vals[0].Int()))
			responseVal = vals[1]
		} else if len(vals) > 0 {
			responseVal = vals[0]
		}
		if canDeref(responseVal) {
			responseVal = responseVal.Elem()
		}
		if isByteSlice(responseVal) {
			res.Write(responseVal.Bytes())
		} else {
			res.Write([]byte(responseVal.String()))
		}
	}
}

func isByteSlice(val reflect.Value) bool {
	return val.Kind() == reflect.Slice && val.Type().Elem().Kind() == reflect.Uint8
}

func canDeref(val reflect.Value) bool {
	return val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr
}

func (c *context) Written() bool {
	return c.rw.Written()
}

func (c *context) run() {

	for c.index <= len(c.handlers) {
		vals, err := c.Invoke(c.handler())
		if err != nil {
			panic(err)
		}
		c.index += 1

		// if the handler returned something, write it to the http response
		if len(vals) > 0 {
			ev := c.Get(reflect.TypeOf(ReturnHandler(nil)))

			handleReturn := ev.Interface().(ReturnHandler)
			handleReturn(c, vals)
		}

		if c.Written() {
			return
		}

	}
}
