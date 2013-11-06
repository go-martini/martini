package martini

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

/* Test Helpers */
func expect(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Errorf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func refute(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		t.Errorf("Did not expect %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func Test_New(t *testing.T) {
	m := New()
	refute(t, m, nil)
}

func Test_Martini_ServeHTTP(t *testing.T) {

	result := ""
	response := httptest.NewRecorder()

	m := New()
	m.Use(func(c Context) {
		result += "foo"
		c.Next()
		result += "ban"
	})
	m.Use(func(c Context) {
		result += "bar"
		c.Next()
		result += "baz"
	})
	m.Use(func(res http.ResponseWriter, req *http.Request) {
		result += "bat"
		res.WriteHeader(400)
	})
	m.ServeHTTP(response, (*http.Request)(nil))

	expect(t, result, "foobarbatbazban")
	expect(t, response.Code, 400)
}
