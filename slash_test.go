package martini

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var newSlashTests = [...][]struct {
	url      string
	code     int
	location string
}{
	{
		{"http://localhost/foo/", http.StatusMovedPermanently, "http://localhost/foo"},
		{"http://localhost/foo", http.StatusOK, ""},
		{"http://localhost/", http.StatusOK, ""},
	},
	{
		{"http://localhost/foo/", http.StatusOK, ""},
		{"http://localhost/foo", http.StatusMovedPermanently, "http://localhost/foo/"},
		{"http://localhost/", http.StatusOK, ""},
	},
}

func Test_NewSlash(t *testing.T) {
	for i, v := range newSlashTests {
		handler := NewSlash(i == 1).(func(http.ResponseWriter, *http.Request))
		for _, tt := range v {
			recorder := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tt.url, nil)
			handler(recorder, req)
			expect(t, recorder.Code, tt.code)
			expect(t, recorder.Header().Get("Location"), tt.location)
		}
	}
}
