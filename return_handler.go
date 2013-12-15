package martini

import (
	"net/http"
	"reflect"
)

// ReturnHandler is a service that Martini provides that is called
// when a route handler returns something. The ReturnHandler is
// responsible for writing to the ResponseWriter based on the values
// that are passed into this function.
type ReturnHandler func(http.ResponseWriter, []reflect.Value)

func defaultReturnHandler() ReturnHandler {
	return func(res http.ResponseWriter, vals []reflect.Value) {
		if len(vals) > 1 && vals[0].Kind() == reflect.Int {
			res.WriteHeader(int(vals[0].Int()))
			res.Write([]byte(vals[1].String()))
		} else if len(vals) > 0 {
			res.Write([]byte(vals[0].String()))
		}
	}
}
