package martini

import (
	"testing"
)

func Test_Routing(t *testing.T) {
	router := NewRouter()
	result := ""
	router.Get("/foo", func() {
		result += "foo"
	})
}
