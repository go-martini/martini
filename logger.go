package martini

import (
	"log"
	"net/http"
	"time"
)

func Logger(res http.ResponseWriter, req *http.Request, c Context, log *log.Logger) {
	t0 := time.Now()
	log.Printf("\033[32;1mStarted %s %s\033[0m\n", req.Method, req.URL.Path)

	rl := &responseLogger{res, 0, 0}
	c.MapTo(rl, (*http.ResponseWriter)(nil))

	c.Next()

	t1 := time.Now()
	log.Printf("Completed %v in %v\n", rl.status, t1.Sub(t0))
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
