package martini

import (
	"net/http"
	"time"
)

// Logger interface compliant with stdlib's log.Logger
// See documentation for log.Logger for more info
type Logger interface {
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})
	Flags() int
	Output(int, string) error
	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
	Prefix() string
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})
	SetFlags(int)
	SetPrefix(string)
}

// LoggerMiddleware returns a middleware handler that logs the request as it goes in and the response as it goes out.
func LoggerMiddleware() Handler {
	return func(res http.ResponseWriter, req *http.Request, c Context, log Logger) {
		start := time.Now()

		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			addr = req.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = req.RemoteAddr
			}
		}

		log.Printf("Started %s %s for %s", req.Method, req.URL.Path, addr)

		rw := res.(ResponseWriter)
		c.Next()

		log.Printf("Completed %v %s in %v\n", rw.Status(), http.StatusText(rw.Status()), time.Since(start))
	}
}
