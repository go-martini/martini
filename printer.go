package martini

import (
	"fmt"
	"io"
	"net/http"
	"reflect"
)

var typeOfValues = reflect.TypeOf(Values{})

// Printer returns a middleware that prints values returned form a handler.
func Printer() Handler {
	return func(res http.ResponseWriter, c Context) {
		c.Next()
		if values, _ := c.Get(typeOfValues).Interface().(Values); len(values) > 0 {
			if status, isInt := values[0].(int); isInt {
				res.WriteHeader(status)
			}
			switch v := values[len(values)-1].(type) {
			case string:
				io.WriteString(res, v)
			case fmt.Stringer:
				io.WriteString(res, v.String())
			case []byte:
				res.Write(v)
			}
		}
	}
}
