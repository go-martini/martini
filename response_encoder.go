package martini

import (
	"net/http"
	"reflect"
)

// ResponseEncoder is a service that Martini provides that is called
// when a route handler returns something. The ResponseEncoder is
// responsible for writing to the ResponseWriter based on the values
// that are passed into this function.
type ResponseEncoder func(http.ResponseWriter, []reflect.Value)

func defaultResponseEncoder(res http.ResponseWriter, vals []reflect.Value) {
	if len(vals) > 1 && vals[0].Kind() == reflect.Int {
		res.WriteHeader(int(vals[0].Int()))
		res.Write([]byte(vals[1].String()))
	} else if len(vals) > 0 {
		res.Write([]byte(vals[0].String()))
	}
}
