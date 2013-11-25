package martini

import (
	"fmt"
	"net/http"
	"reflect"
	"regexp"
)

// Params is a map of name/value pairs for named routes. An instance of martini.Params is available to be injected into any route handler.
type Params map[string]string

// Router is Martini's de-facto routing interface. Supports HTTP verbs, stacked handlers, and dependency injection.
type Router interface {
	// Get adds a route for a HTTP GET request to the specified matching pattern.
	Get(string, ...Handler) Route
	// Patch adds a route for a HTTP PATCH request to the specified matching pattern.
	Patch(string, ...Handler) Route
	// Post adds a route for a HTTP POST request to the specified matching pattern.
	Post(string, ...Handler) Route
	// Put adds a route for a HTTP PUT request to the specified matching pattern.
	Put(string, ...Handler) Route
	// Delete adds a route for a HTTP DELETE request to the specified matching pattern.
	Delete(string, ...Handler) Route

	// NotFound sets the handler that is called when a no route matches a request. Throws a basic 404 by default.
	NotFound(Handler)

	// Handle is the entry point for routing. This is used as a martini.Handler
	Handle(http.ResponseWriter, *http.Request, Context)
}

type router struct {
	routes   []Route
	notFound Handler
}

// NewRouter creates a new Router instance.
func NewRouter() Router {
	return &router{notFound: http.NotFound}
}

func (r *router) Get(pattern string, h ...Handler) Route {
	return r.addRoute("GET", pattern, h)
}

func (r *router) Patch(pattern string, h ...Handler) Route {
	return r.addRoute("PATCH", pattern, h)
}

func (r *router) Post(pattern string, h ...Handler) Route {
	return r.addRoute("POST", pattern, h)
}

func (r *router) Put(pattern string, h ...Handler) Route {
	return r.addRoute("PUT", pattern, h)
}

func (r *router) Delete(pattern string, h ...Handler) Route {
	return r.addRoute("DELETE", pattern, h)
}

func (r *router) Handle(res http.ResponseWriter, req *http.Request, context Context) {
	for _, route := range r.routes {
		ok, vals := route.Match(req.Method, req.URL.Path)
		if ok {
			params := Params(vals)
			context.Map(params)
            rh := Routes{}
            context.Map(rh)
			_, err := context.Invoke(route.Handle)
			if err != nil {
				panic(err)
			}
			return
		}
	}

	// no routes exist, 404
	_, err := context.Invoke(r.notFound)
	if err != nil {
		panic(err)
	}
}

func (r *router) NotFound(handler Handler) {
	r.notFound = handler
}

func (r *router) addRoute(method string, pattern string, handlers []Handler) *route {
	route := newRoute(method, pattern, handlers)
	route.Validate()
	r.routes = append(r.routes, route)
	return route
}

// Route is the default route interface.
type Route interface {
    UrlWith([]string) string
    Match(string, string) (bool, map[string]string)
    Validate()
    Handle(Context, http.ResponseWriter)
}

type route struct {
	method    string
	regex     *regexp.Regexp
	handlers  []Handler
	RouteName string
	Pattern   string
}

func newRoute(method string, pattern string, handlers []Handler) *route {
	route := route{method, nil, handlers, "", pattern}
	r := regexp.MustCompile(`:[^/#?()\.\\]+`)
	pattern = r.ReplaceAllStringFunc(pattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:len(m)])

	})
	pattern += `\/?`
	route.regex = regexp.MustCompile(pattern)
	return &route
}

func (r *route) Match(method string, path string) (bool, map[string]string) {
	if method != r.method {
		return false, nil
	}

	matches := r.regex.FindStringSubmatch(path)
	if len(matches) > 0 && matches[0] == path {
		params := make(map[string]string)
		for i, name := range r.regex.SubexpNames() {
			if len(name) > 0 {
				params[name] = matches[i]
			}
		}
		return true, params
	}
	return false, nil
}

func (r *route) Validate() {
	for _, handler := range r.handlers {
		validateHandler(handler)
	}
}

func (r *route) Handle(c Context, res http.ResponseWriter) {
	for _, handler := range r.handlers {
		vals, err := c.Invoke(handler)
		if err != nil {
			panic(err)
		}

		// if the handler returned something, write it to
		// the http response
		if len(vals) > 1 && vals[0].Kind() == reflect.Int {
			res.WriteHeader(int(vals[0].Int()))
			res.Write([]byte(vals[1].String()))
		} else if len(vals) > 0 {
			res.Write([]byte(vals[0].String()))
		}
		if c.written() {
			return
		}
	}
}

// UrlWith returns the url pattern replacing the parameters for its values
func (r *route) UrlWith(args []string) string {
	if len(args) > 0 {
		reg := regexp.MustCompile(`:[^/#?()\.\\]+`)
		argCount := len(args)
		i := 0
		url := reg.ReplaceAllStringFunc(r.Pattern, func(m string) string {
			var val interface{}
			if i < argCount {
				val = args[i]
			} else {
				val = m
			}
			i += 1
			return fmt.Sprintf(`%v`, val)
		})

		return url
	} else {
		return r.Pattern
	}
}
