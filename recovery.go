package martini

import (
	"net/http"
)

func RecoveryHandler() Handler {
	return func(res http.ResponseWriter, c Context) {
		defer func() {
			if err := recover(); err != nil {
				res.WriteHeader(500)
			}
		}()

		c.Next()
	}
}
