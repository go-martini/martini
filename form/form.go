// Package form implements a handler for post and get forms.
//
// For a full guide visit http://github.com/codegangsta/martini
//
//  package main
//
//  import (
//     "github.com/codegangsta/martini"
//     "github.com/codegangsta/martini/form"
//   )
//
//  type BlogPost struct{
//     Title string `form:"title,required"`
//     Content string `form:"content"`
//  }
//
//  func main() {
//    m := martini.Classic()
//
//    m.Get("/",form.Form(&BlogPost{}) func(blogpost *BlogPost) string {
//      return blogpost.Title
//    })
//
//    m.Run()
//  }
package form

import (
	"github.com/codegangsta/martini"
	"net/http"
	"reflect"
	"strings"
)

const (
	RequireError string = "RequireError" //Error for fields which are marked as required but weren't present in the form.
)

//Available errors. Use len() to check if any errors occured.
type Errors map[string]string

//Create a new formhandler. Errors are available via form.Errors-Service.
func Form(formstruct interface{}) martini.Handler {
	return func(context martini.Context, req *http.Request) {
		req.ParseForm()
		typ := reflect.TypeOf(formstruct).Elem()
		errors := make(Errors)

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if tag := field.Tag.Get("form"); tag != "" {
				args := strings.Split(tag, ",")
				if len(args) > 0 {
					val := req.Form.Get(args[0])
					reflect.ValueOf(formstruct).Elem().Field(i).SetString(val)
					if len(args) > 1 {
						if val == "" && strings.Contains(args[1], "required") {
							errors[args[0]] = RequireError
						}
					}
				}
			}
		}
		context.Map(errors)
		context.Map(formstruct)
	}
}
