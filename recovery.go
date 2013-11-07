package martini

import (
	"log"
	"net/http"
)

func Recovery() Handler {
	return func(res http.ResponseWriter, c Context, logger *log.Logger) {
		defer func() {
			if err := recover(); err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				logger.Printf("PANIC: %v\n", err)
			}
		}()

		c.Next()
	}
}
