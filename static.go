package martini

import (
	"log"
	"net/http"
	"path"
	"strings"
)

// StaticOptions is a struct for specifying configuration options for the martini.Static middleware.
type StaticOptions struct {
	// SkipLogging can be used to switch log messages to *log.logger off.
	SkipLogging bool
	// SkipServeIndex is used to determine whether or not an index file should be tried to serve by the martini.Static middleware.
	SkipServeIndex bool
	// IndexFile defines which file to serve as index if it exists. Has no effect if SkipServeIndex is true.
	IndexFile string
}

func prepareStaticOptions(options []StaticOptions) StaticOptions {
	var opt StaticOptions
	if len(options) > 0 {
		opt = options[0]
	}

	// Defaults
	if len(opt.IndexFile) == 0 {
		opt.IndexFile = "index.html"
	}

	return opt
}

// Static returns a middleware handler that serves static files in the given directory.
func Static(directory string, staticOpt ...StaticOptions) Handler {
	dir := http.Dir(directory)
	opt := prepareStaticOptions(staticOpt)

	return func(res http.ResponseWriter, req *http.Request, log *log.Logger) {
		if req.Method != "GET" && req.Method != "HEAD" {
			return
		}
		// Do nothing if we are not serving any index file and the request is at the root URL path.
		if opt.SkipServeIndex && req.URL.Path == "/" {
			return
		}
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

		if !opt.SkipServeIndex {
			// try to serve index file
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
