package martini

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// Static returns a middleware handler that serves static files in the given path.
func Static(path string) Handler {
	return func(res http.ResponseWriter, req *http.Request, log *log.Logger) {
		file := filepath.Join(path, filepath.Clean(req.URL.Path))
		info, err := os.Stat(file)
		if err == nil && !info.IsDir() {
			log.Println("[Static] Serving " + file)
			http.ServeFile(res, req, file)
		}
	}
}
