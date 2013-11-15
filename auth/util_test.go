package auth

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func Test_SecureCompare(t *testing.T) {
	Convey("SecureCompare should correctly assert value equality", t, func() {
		for _, tt := range comparetests {
			So(SecureCompare(tt.a, tt.b), ShouldEqual, tt.val)
		}
	})
}

var comparetests = []struct {
	a   string
	b   string
	val bool
}{
	{"foo", "foo", true},
	{"bar", "bar", true},
	{"password", "password", true},
	{"Foo", "foo", false},
	{"foo", "foobar", false},
	{"password", "pass", false},
}
