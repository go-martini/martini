package martini

import (
	"net/http"
)

type Router interface {
	Get(string, Handler)
	Post(string, Handler)
	Put(string, Handler)
	Delete(string, Handler)

	Handle(Context)
}

type route struct {
	method  string
	pattern string
	handler Handler
}

type router struct {
	routes []route
}

func (r *router) Get(pattern string, handler Handler) {
	r.addRoute("GET", pattern, handler)
}

func (r *router) Post(pattern string, handler Handler) {
	r.addRoute("POST", pattern, handler)
}

func (r *router) Put(pattern string, handler Handler) {
	r.addRoute("PUT", pattern, handler)
}

func (r *router) Delete(pattern string, handler Handler) {
	r.addRoute("DELETE", pattern, handler)
}

func (r *router) Handle(context Context, req *http.Request) {
	for _, route := range r.routes {
		// Be super strict for now. Eventually we will have some
		// super awesome pattern matching here. But not today
		if route.method == req.Method && req.URL.Path == route.pattern {
			err := context.Invoke(route.handler)
			if err != nil {
				panic(err)
			}
			return
		}
	}
}

func (r *router) addRoute(method string, pattern string, handler Handler) {
	r.routes = append(r.routes, route{method, pattern, handler})
}
