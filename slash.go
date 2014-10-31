package martini

import "net/http"

// NewSlash return a middleware to take care of last slash.
// When add is true always add last slash, false remove it.
func NewSlash(add bool) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "CONNECT" {
			u := r.URL
			if n := len(u.Path); u.Path != "/" && n > 0 {
				n--
				if u.Path[n] == '/' {
					if add {
						return
					} else {
						u.Path = u.Path[:n]
					}
				} else if add {
					u.Path += "/"
				} else {
					return
				}
				http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
			}
		}
	}
}
