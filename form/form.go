//Package form implements a handler for post and get forms.
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
	"fmt"
	"martini"
	"net/http"
	"reflect"
	"strings"
)

//Avaiable error if a formfield which was set as required, wasn't present.
type RequireError struct {
	Fields []string //Name of fields which weren't set
}

func newRequireError(fields []string) *RequireError {
	return &RequireError{
		fields,
	}
}

func (re *RequireError) Error() string {
	return fmt.Sprintf("Fields: %v are required", re.Fields)
}

//Create a new formhandler. If a required field isn't present, a *RequireError is avaiable.
func Form(formstruct interface{}) martini.Handler {
	return func(context martini.Context, req *http.Request) {
		req.ParseForm()
		typ := reflect.TypeOf(formstruct).Elem()
		errfields := make([]string, 0)

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			if tag := field.Tag.Get("form"); tag != "" {
				args := strings.Split(tag, ",")
				if len(args) > 0 {
					val := req.Form.Get(args[0])
					reflect.ValueOf(formstruct).Elem().Field(i).SetString(val)
					if len(args) > 1 {
						if val == "" && strings.Contains(args[1], "required") {
							errfields = append(errfields, args[0])
						}
					}
				}
			}

		}

		if len(errfields) > 0 {
			context.Map(newRequireError(errfields))
		} else {
			context.Map((*RequireError)(nil))
		}
		context.Map(formstruct)
	}
}
