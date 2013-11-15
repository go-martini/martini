package form

import (
	"martini"
	"net/http"
	"net/http/httptest"
	"testing"
)

type formTest struct {
	method string
	path   string
	ok     bool
	ref    *BlogPost
}

var formTests = []formTest{
	{"GET", "http://localhost:3000/blogposts/create?content=Test", false, &BlogPost{"", "Test"}},
	{"POST", "http://localhost:3000/blogposts/create?content=Test&title=TheTitle", true, &BlogPost{"TheTitle", "Test"}},
}

type BlogPost struct {
	Title   string `form:"title,required"`
	Content string `form:"content"`
}

func assertEqualField(t *testing.T, fieldname string, testcasenumber int, expected interface{}, got interface{}) {
	if expected != got {
		t.Errorf("%s: expected=%s, got=%s in Testcase:%i\n", fieldname, expected, got, testcasenumber)
	}
}

func handler(test formTest, t *testing.T, index int, post *BlogPost, errors Errors) {
	if !test.ok && len(errors) == 0 {
		t.Errorf("expected RequireError in Testcase:%i", index)
	}
	assertEqualField(t, "Title", index, test.ref.Title, post.Title)
	assertEqualField(t, "Content", index, test.ref.Content, post.Content)

}

func Test_FormTests(t *testing.T) {
	for index, test := range formTests {
		recorder := httptest.NewRecorder()
		m := martini.Classic()
		if test.method == "GET" {
			m.Get("/blogposts/create", Form(&BlogPost{}), func(post *BlogPost, errors Errors) { handler(test, t, index, post, errors) })
		}
		if test.method == "POST" {
			m.Post("/blogposts/create", Form(&BlogPost{}), func(post *BlogPost, errors Errors) { handler(test, t, index, post, errors) })
		}
		req, err := http.NewRequest(test.method, test.path, nil)
		if err != nil {
			t.Error(err)
		}
		m.ServeHTTP(recorder, req)
	}
}
