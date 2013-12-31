package martini

import (
	"fmt"
	"github.com/codegangsta/inject"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
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
	// Options adds a route for a HTTP OPTIONS request to the specified matching pattern.
	Options(string, ...Handler) Route
	// Head adds a route for a HTTP HEAD request to the specified matching pattern.
	Head(string, ...Handler) Route
	// Any adds a route for any HTTP method request to the specified matching pattern.
	Any(string, ...Handler) Route

	// NotFound sets the handlers that are called when a no route matches a request. Throws a basic 404 by default.
	NotFound(...Handler)

	// Handle is the entry point for routing. This is used as a martini.Handler
	Handle(http.ResponseWriter, *http.Request, Context)
}

type router struct {
	routes    []*route
	notFounds []Handler
}

// NewRouter creates a new Router instance.
func NewRouter() Router {
	return &router{notFounds: []Handler{http.NotFound}}
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

func (r *router) Options(pattern string, h ...Handler) Route {
	return r.addRoute("OPTIONS", pattern, h)
}

func (r *router) Head(pattern string, h ...Handler) Route {
	return r.addRoute("HEAD", pattern, h)
}

func (r *router) Any(pattern string, h ...Handler) Route {
	return r.addRoute("*", pattern, h)
}

func (r *router) Handle(res http.ResponseWriter, req *http.Request, context Context) {
	for _, route := range r.routes {
		ok, vals := route.Match(req.Method, req.URL.Path)
		if ok {
			params := Params(vals)
			context.Map(params)
			r := routes{}
			context.MapTo(r, (*Routes)(nil))
			_, err := context.Invoke(route.Handle)
			if err != nil {
				panic(err)
			}
			return
		}
	}

	// no routes exist, 404
	c := &routeContext{context, 0, r.notFounds}
	context.MapTo(c, (*Context)(nil))
	c.run()
}

func (r *router) NotFound(handler ...Handler) {
	r.notFounds = handler
}

func (r *router) addRoute(method string, pattern string, handlers []Handler) *route {
	route := newRoute(method, pattern, handlers)
	route.Validate()
	r.routes = append(r.routes, route)
	return route
}

// Route is an interface representing a Route in Martini's routing layer.
type Route interface {
	// URLWith returns a rendering of the Route's url with the given string params.
	URLWith([]string) string
}

type route struct {
	method   string
	regex    *regexp.Regexp
	handlers []Handler
	pattern  string
}

func newRoute(method string, pattern string, handlers []Handler) *route {
	route := route{method, nil, handlers, pattern}
	r := regexp.MustCompile(`:[^/#?()\.\\]+`)
	pattern = r.ReplaceAllStringFunc(pattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:])
	})
	r2 := regexp.MustCompile(`\*\*`)
	var index int
	pattern = r2.ReplaceAllStringFunc(pattern, func(m string) string {
		index++
		return fmt.Sprintf(`(?P<_%d>[^#?]*)`, index)
	})
	pattern += `\/?`
	route.regex = regexp.MustCompile(pattern)
	return &route
}

func (r route) MatchMethod(method string) bool {
	return r.method == "*" || method == r.method || (method == "HEAD" && r.method == "GET")
}

func (r route) Match(method string, path string) (bool, map[string]string) {
	// add Any method matching support
	if !r.MatchMethod(method) {
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
	context := &routeContext{c, 0, r.handlers}
	c.MapTo(context, (*Context)(nil))
	context.run()
}

// URLWith returns the url pattern replacing the parameters for its values
func (r *route) URLWith(args []string) string {
	if len(args) > 0 {
		reg := regexp.MustCompile(`:[^/#?()\.\\]+`)
		argCount := len(args)
		i := 0
		url := reg.ReplaceAllStringFunc(r.pattern, func(m string) string {
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
	}
	return r.pattern
}

// Routes is a helper service for Martini's routing layer.
type Routes interface {
	// URLFor returns a rendered URL for the given route. Optional params can be passed to fulfill named parameters in the route.
	URLFor(route Route, params ...interface{}) string
}

type routes struct{}

// URLFor returns the url for the given route name.
func (r routes) URLFor(route Route, params ...interface{}) string {
	var args []string
	for _, param := range params {
		switch v := param.(type) {
		case int:
			args = append(args, strconv.FormatInt(int64(v), 10))
		case string:
			args = append(args, v)
		default:
			if v != nil {
				panic("Arguments passed to UrlFor must be integers or strings")
			}
		}
	}

	return route.URLWith(args)
}

type routeContext struct {
	Context
	index    int
	handlers []Handler
}

func (r *routeContext) Next() {
	r.index += 1
	r.run()
}

func (r *routeContext) run() {
	for r.index < len(r.handlers) {
		handler := r.handlers[r.index]
		vals, err := r.Invoke(handler)
		if err != nil {
			panic(err)
		}
		r.index += 1

		// if the handler returned something, write it to the http response
		if len(vals) > 0 {
			rv := r.Get(inject.InterfaceOf((*http.ResponseWriter)(nil)))
			ev := r.Get(reflect.TypeOf(ReturnHandler(nil)))
			handleReturn := ev.Interface().(ReturnHandler)
			handleReturn(rv.Interface().(http.ResponseWriter), vals)
		}

		if r.Written() {
			return
		}
	}
}
