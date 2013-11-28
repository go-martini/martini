package martini

import (
	"log"
	"net/http"
	"path/filepath"
)

// Static returns a middleware handler that serves static files in the given path.
func Static(path string) Handler {
	dir := http.Dir(path)
	return func(res http.ResponseWriter, req *http.Request, log *log.Logger) {
		file := req.URL.Path
		f, err := dir.Open(file)
		if err != nil {
			// discard the error?
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			return
		}

		// Try to serve index.html
		if fi.IsDir() {

			// redirect if missing trailing slash
			if file[len(file)-1] != '/' {
				http.Redirect(res, req, file+"/", http.StatusFound)
				return
			}

			file = filepath.Join(file, "index.html")
			f, err = dir.Open(file)
			if err != nil {
				return
			}
			defer f.Close()

			fi, err = f.Stat()
			if err != nil || fi.IsDir() {
				return
			}
		}

		log.Println("[Static] Serving " + file)
		http.ServeContent(res, req, file, fi.ModTime(), f)
	}
}
