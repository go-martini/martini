package martini

import "net/http"

// When x is true always add last slash, false remove.
func NewSlash(x bool) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "CONNECT" {
			u := r.URL
			if n := len(u.Path); u.Path != "/" && n > 0 {
				n--
				if u.Path[n] == '/' {
					if !x {
						u.Path = u.Path[:n]
					} else {
						return
					}
				} else if x {
					u.Path += "/"
				} else {
					return
				}
				http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
			}
		}
	}
}
