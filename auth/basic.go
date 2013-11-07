package auth

import (
	"net/http"
)

func Basic(name string, password string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		u := req.URL.User.Username()
		p, _ := req.URL.User.Password()

		if u != name || p != password {
			res.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
			res.WriteHeader(401)
		}
	}
}
