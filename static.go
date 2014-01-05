package martini

import (
	"log"
	"net/http"
	"path"
	"strings"
)

type StaticOptions struct {
	SkipLogging    bool
	SkipServeIndex bool
	IndexFile      string
}

func prepareStaticOptions(options []StaticOptions) StaticOptions {
	if len(options) > 0 {
		// Default index file to serve should be index.html
		if options[0].IndexFile == "" {
			options[0].IndexFile = "index.html"
		}
		return options[0]
	}
	return StaticOptions{
		SkipLogging:    false,        // Logging is on by default
		SkipServeIndex: false,        // Try to serve index.html by default
		IndexFile:      "index.html", // Default index file to serve
	}
}

// Static returns a middleware handler that serves static files in the given directory.
func Static(directory string, staticOpt ...StaticOptions) Handler {
	dir := http.Dir(directory)
	opt := prepareStaticOptions(staticOpt)

	return func(res http.ResponseWriter, req *http.Request, log *log.Logger) {
		file := req.URL.Path
		f, err := dir.Open(file)
		if err != nil || req.Method != "GET" {
			// discard the error?
			return
		}
		defer f.Close()

		fi, err := f.Stat()
		if err != nil {
			return
		}

		if !opt.SkipServeIndex {
			// Try to serve index.html
			if fi.IsDir() {

				// redirect if missing trailing slash
				if !strings.HasSuffix(file, "/") {
					http.Redirect(res, req, file+"/", http.StatusFound)
					return
				}

				file = path.Join(file, opt.IndexFile)
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
		}

		if !opt.SkipLogging {
			log.Println("[Static] Serving " + file)
		}
		http.ServeContent(res, req, file, fi.ModTime(), f)
	}
}
