package martini

import (
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

func Test_NewApp(t *testing.T) {
	m := New()
	refute(t, m, nil)
}

func Test_App_Use(t *testing.T) {
	handleFunc := func() {
	}

	m := New()
	m.Use(handleFunc)
	expect(t, len(m.handlers), 1)
}
