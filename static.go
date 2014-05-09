package martini

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
	"time"
)

// StaticOptions is a struct for specifying configuration options for the martini.Static middleware.
type StaticOptions struct {
	// Prefix is the optional prefix used to serve the static directory content
	Prefix string
	// SkipLogging will disable [Static] log messages when a static file is served.
	SkipLogging bool
	// IndexFile defines which file to serve as index if it exists.
	IndexFile string
	// Expires defines which user-defined function to use for producing a HTTP Expires Header
	// https://developers.google.com/speed/docs/insights/LeverageBrowserCaching
	Expires func() string
	// BinData is a map of a path and its content (binary data).
	// If this field is set, Static tries to retrieve data from this field on memory
	// instead of files on disk.
	BinData map[string][]byte
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
	// Normalize the prefix if provided
	if opt.Prefix != "" {
		// Ensure we have a leading '/'
		if opt.Prefix[0] != '/' {
			opt.Prefix = "/" + opt.Prefix
		}
		// Remove any trailing '/'
		opt.Prefix = strings.TrimRight(opt.Prefix, "/")
	}
	return opt
}

// Static returns a middleware handler that serves static files in the given directory.
func Static(directory string, staticOpt ...StaticOptions) Handler {
	dir := http.Dir(directory)
	opt := prepareStaticOptions(staticOpt)

	// set current time to binDataModtime if opt.BinData is set.
	var binDataModtime time.Time
	if opt.BinData != nil {
		binDataModtime = time.Now()
	}

	return func(res http.ResponseWriter, req *http.Request, log *log.Logger) {
		if req.Method != "GET" && req.Method != "HEAD" {
			return
		}
		file := req.URL.Path
		// if we have a prefix, filter requests by stripping the prefix
		if opt.Prefix != "" {
			if !strings.HasPrefix(file, opt.Prefix) {
				return
			}
			file = file[len(opt.Prefix):]
			if file != "" && file[0] != '/' {
				return
			}
		}

		// if binary data of the file exists, serve this data.
		if opt.BinData != nil && len(opt.BinData[file]) != 0 {
			serveContent(res, req, file, binDataModtime, bytes.NewReader(opt.BinData[file]), opt, log)
			return
		}

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

		// try to serve index file
		if fi.IsDir() {
			// redirect if missing trailing slash
			if !strings.HasSuffix(req.URL.Path, "/") {
				http.Redirect(res, req, req.URL.Path+"/", http.StatusFound)
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

		serveContent(res, req, file, fi.ModTime(), f, opt, log)
	}
}

// serveContent calls http.ServeContent.
func serveContent(
	res http.ResponseWriter,
	req *http.Request,
	name string,
	modtime time.Time,
	content io.ReadSeeker,
	opt StaticOptions,
	log *log.Logger,
) {
	if !opt.SkipLogging {
		log.Println("[Static] Serving " + name)
	}

	// Add an Expires header to the static content
	if opt.Expires != nil {
		res.Header().Set("Expires", opt.Expires())
	}

	http.ServeContent(res, req, name, modtime, content)
}
