package martini

import (
	"fmt"
	"net/http"
	"regexp"
)

type Router interface {
	Get(string, ...Handler)
	Post(string, ...Handler)
	Put(string, ...Handler)
	Delete(string, ...Handler)

	Handle(http.ResponseWriter, *http.Request, Context)
}

type router struct {
	routes []route
}

func NewRouter() Router {
	return &router{}
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
		ok, _ := route.match(req.Method, req.URL.Path)
		if ok {
			_, err := context.Invoke(route.handle)
			if err != nil {
				panic(err)
			}
			return
		}
	}

	// no routes exist, 404
	res.WriteHeader(http.StatusNotFound)
}

func (r *router) addRoute(method string, pattern string, handlers []Handler) {
	route := newRoute(method, pattern, handlers)
	if route.validate() == nil {
		r.routes = append(r.routes, route)
	}
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

func (r route) validate() error {
	for _, handler := range r.handlers {
		err := validateHandler(handler)
		if err != nil {
			return err
		}
	}
	return nil
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
