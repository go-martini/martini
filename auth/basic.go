package auth

import (
	"encoding/base64"
	"net/http"
	"strings"
)

func Basic(name string, password string) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		auth, ok := req.Header["Authorization"]
		if !ok {
			unauthorized(res)
		}

		u, p := parseToken(auth[0])

		if u != name || p != password {
			unauthorized(res)
		}
	}
}

func unauthorized(res http.ResponseWriter) {
	res.Header().Set("WWW-Authenticate", "Basic realm=\"Authorization Required\"")
	res.WriteHeader(http.StatusUnauthorized)
}

// ParseToken is a helper function that extracts the username and password
// from an authorization token.  Callers should be able to pass in the header
// "Authorization" from an HTTP request, and retrieve the credentials.
func parseToken(token string) (username, password string) {
	if token == "" {
		return "", ""

	}

	// Check that the token supplied corresponds to the basic authorization
	// protocol
	ndx := strings.IndexRune(token, ' ')
	if ndx < 1 || token[0:ndx] != "Basic" {
		return "", ""

	}

	// Drop prefix, and decode the base64
	buffer, err := base64.StdEncoding.DecodeString(token[ndx+1:])
	if err != nil {
		return "", ""

	}
	token = string(buffer)

	ndx = strings.IndexRune(token, ':')
	if ndx < 1 {
		return "", ""

	}

	return token[0:ndx], token[ndx+1:]

}
