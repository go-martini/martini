package martini

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
)

const (
	panicHtml = `<html>
<head><title>PANIC: %s</title></head>
<style type="text/css">
html, body {
	font-family: "Roboto", sans-serif;
	color: #333333;
	background-color: #ea5343;
	margin: 0px;
}
h1 {
	color: #d04526;
	background-color: #ffffff;
	padding: 20px;
	border-bottom: 1px dashed #2b3848;
}
pre {
	margin: 20px;
	padding: 20px;
	border: 2px solid #2b3848;
	background-color: #ffffff;
}
</style>
<body>
<h1>PANIC</h1>
<pre style="font-weight: bold;">%s</pre>
<pre>%s</pre>
</body>
</html>`
)

// Recovery returns a middleware that recovers from any panics and writes a 500 if there was one.
// While Martini is in development mode, Recovery will also output the panic as HTML.
func Recovery() Handler {
	return func(res http.ResponseWriter, c Context, log *log.Logger) {
		defer func() {
			if err := recover(); err != nil {
				res.WriteHeader(http.StatusInternalServerError)
				log.Printf("PANIC: %s\n%s", err, debug.Stack())

				// respond with panic message while in development mode
				if Env == Dev {
					res.Write([]byte(fmt.Sprintf(panicHtml, err, err, debug.Stack())))
				}
			}
		}()

		c.Next()
	}
}
