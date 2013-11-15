package martini

import (
	"log"
	"net/http"
	"runtime"
	"time"
)

var logFormat = getLogFormat()

func getLogFormat() string {
	format := "Started %s %s"
	if runtime.GOOS != "windows" {
		return "\033[32;1m" + format + "\033[0m\n"
	}
	return format
}

// Logger returns a middleware handler that logs the request as it goes in and the response as it goes out.
func Logger() Handler {
	return func(res http.ResponseWriter, req *http.Request, c Context, log *log.Logger) {
		start := time.Now()
		log.Printf(logFormat, req.Method, req.URL.Path)

		rl := &responseLogger{res, 200, 0}
		c.MapTo(rl, (*http.ResponseWriter)(nil))

		c.Next()

		log.Printf("Completed %v %s in %v\n", rl.status, http.StatusText(rl.status), time.Now().Sub(start))
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
