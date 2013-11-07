package auth

import (
	"net/http"
)

func Basic(name string, password string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		u := req.URL.User.Username()
		p, _ := req.URL.User.Password()

		if u != name || p != password {
			res.WriteHeader(401)
		}
	}
}
