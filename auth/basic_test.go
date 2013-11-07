package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_BasicAuth(t *testing.T) {
	res := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "http://hello:world@localhost:3000/foobar", nil)
	if err != nil {
		t.Error(err)
	}

	b := Basic("hello", "world")
	b(res, req)

	if res.Code != 200 {
		t.Errorf("response code is not 0 : %v", res.Code)
	}

	b2 := Basic("", "")
	b2(res, req)

	if res.Code != 401 {
		t.Errorf("response code is not 403 : %v", res.Code)
	}

	if res.Header().Get("WWW-Authenticate") != "Basic realm=\"Authorization Required\"" {
		t.Error("Authentication header missing")
	}
}
