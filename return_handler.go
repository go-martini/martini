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
		var responseVal reflect.Value
		if len(vals) > 1 && vals[0].Kind() == reflect.Int {
			res.WriteHeader(int(vals[0].Int()))
			responseVal = vals[1]
		} else if len(vals) > 0 {
			responseVal = vals[0]
		}
		if responseVal.Kind() == reflect.Interface || responseVal.Kind() == reflect.Ptr {
			responseVal = responseVal.Elem()
		}
		if responseVal.Kind() == reflect.Slice && responseVal.Type().Elem().Kind() == reflect.Uint8 {
			res.Write(responseVal.Bytes())
		} else {
			res.Write([]byte(responseVal.String()))
		}
	}
}
