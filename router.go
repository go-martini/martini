package martini

import (
	"fmt"
	"net/http"
	"regexp"
)

// Params is a map of name/value pairs for named routes. An instance of martini.Params is available to be injected into any route handler.
type Params map[string]string

// Router is Martini's de-facto routing interface. Supports HTTP verbs, stacked handlers, and dependency injection.
type Router interface {
	// Get adds a route for a HTTP GET request to the specified matching pattern.
	Get(string, ...Handler)
	// Post adds a route for a HTTP POST request to the specified matching pattern.
	Post(string, ...Handler)
	// Put adds a route for a HTTP PUT request to the specified matching pattern.
	Put(string, ...Handler)
	// Delete adds a route for a HTTP DELETE request to the specified matching pattern.
	Delete(string, ...Handler)

	// NotFound sets the handler that is called when a no route matches a request. Throws a basic 404 by default.
	NotFound(Handler)

	// Handle is the entry point for routing. This is used as a martini.Handler
	Handle(http.ResponseWriter, *http.Request, Context)
}

type router struct {
	routes   []route
	notFound Handler
}

// NewRouter creates a new Router instance.
func NewRouter() Router {
	return &router{notFound: http.NotFound}
}

func (r *router) Get(pattern string, h ...Handler) {
	r.addRoute("GET", pattern, h)
}

func (r *router) Post(pattern string, h ...Handler) {
	r.addRoute("POST", pattern, h)
}

func (r *router) Put(pattern string, h ...Handler) {
	r.addRoute("PUT", pattern, h)
}

func (r *router) Delete(pattern string, h ...Handler) {
	r.addRoute("DELETE", pattern, h)
}

func (r *router) Handle(res http.ResponseWriter, req *http.Request, context Context) {
	for _, route := range r.routes {
		ok, vals := route.match(req.Method, req.URL.Path)
		if ok {
			params := Params(vals)
			context.Map(params)
			_, err := context.Invoke(route.handle)
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

func (r *router) addRoute(method string, pattern string, handlers []Handler) {
	route := newRoute(method, pattern, handlers)
	route.validate()
	r.routes = append(r.routes, route)
}

type route struct {
	method   string
	regex    *regexp.Regexp
	handlers []Handler
}

func newRoute(method string, pattern string, handlers []Handler) route {
	route := route{method, nil, handlers}
	r := regexp.MustCompile(`:[^/#?()\.\\]+`)
	pattern = r.ReplaceAllStringFunc(pattern, func(m string) string {
		return fmt.Sprintf(`(?P<%s>[^/#?]+)`, m[1:len(m)])

	})
	pattern += `\/?`
	route.regex = regexp.MustCompile(pattern)
	return route
}

func (r route) match(method string, path string) (bool, map[string]string) {
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

func (r route) validate() {
	for _, handler := range r.handlers {
		validateHandler(handler)
	}
}

func (r route) handle(c Context, res http.ResponseWriter) {
	for _, handler := range r.handlers {
		vals, err := c.Invoke(handler)
		if err != nil {
			panic(err)
		}

		// if the handler returned something, write it to
		// the http response
		if len(vals) > 0 {
			res.Write([]byte(vals[0].String()))
		}
		if c.written() {
			return
		}
	}
}
