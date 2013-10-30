package recovery

import (
	"github.com/codegangsta/martini"
	"net/http"
)

func New() martini.Handler {
	return func(res http.ResponseWriter, c martini.Context) {
		defer handlePanic(res)
		c.Next()
	}
}

func handlePanic(res http.ResponseWriter) {
	if err := recover(); err != nil {
		res.WriteHeader(500)
	}
}
