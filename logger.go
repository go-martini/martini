package martini

import (
	"log"
	"net/http"
)

func LoggingHandler() Handler {
	return func(res http.ResponseWriter, req *http.Request, c Context, log *log.Logger) {
		// log the request
		log.Printf("\033[32;1m%s %s\033[0m\n", req.Method, req.URL.Path)

		// override the response writer with a wrapped one
		rl := &responseLogger{res, 0, 0}
		c.MapTo(rl, (*http.ResponseWriter)(nil))

		c.Next()

		// log from responseLogger
		log.Printf("%v - %v bytes \n\t\t%v\n", rl.status, rl.size, rl.Header())
	}
}

type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	if l.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		l.status = http.StatusOK

	}
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}
