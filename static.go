package martini

import (
	"io"
	"log"
	"net/http"
	"path"
	"strings"
)

// Static returns a middleware handler that serves static files in the given directory.
func Static(directory string) Handler {
	dir := http.Dir(directory)

	closeOrPanic := func(c io.Closer) {
		if err := c.Close(); err != nil {
			panic(err)
		}
	}

	return func(res http.ResponseWriter, req *http.Request, log *log.Logger) {
		file := req.URL.Path
		f, err := dir.Open(file)
		if err != nil {
			// discard the error?
			return
		}
		defer closeOrPanic(f)

		fi, err := f.Stat()
		if err != nil {
			return
		}

		// Try to serve index.html
		if fi.IsDir() {

			// redirect if missing trailing slash
			if !strings.HasSuffix(file, "/") {
				http.Redirect(res, req, file+"/", http.StatusFound)
				return
			}

			file = path.Join(file, "index.html")
			f, err = dir.Open(file)
			if err != nil {
				return
			}
			defer closeOrPanic(f)

			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				return
			}
		}

		log.Println("[Static] Serving " + file)
		http.ServeContent(res, req, file, fi.ModTime(), f)
	}
}
