# Martini [![Build Status](https://drone.io/github.com/codegangsta/martini/status.png)](https://drone.io/github.com/codegangsta/martini/latest)

Martini is a package for quickly writing modular web applications/services in Golang.

~~~ go
package main

import "github.com/codegangsta/martini"

func main() {
  m := martini.Classic()
  m.Get("/", func() string {
    return "Hello world!"
  })
  m.Run()
}
~~~

Install the package:
~~~
go get github.com/codegangsta/martini
~~~

## Table of Contents
* [Martini](#martini-)
  * [Table of Contents](#table-of-contents)
  * [Classic Martini](#classic-martini)
    * [Handlers](#handlers)
    * [Routing](#routing)
    * [Services](#services)
    * [Serving Static Files](#serving-static-files)
  * [Middleware Handlers](#middleware-handlers)
    * [Next()](#next)
    * [Injecting Services](#injecting-services)

## Classic Martini
To get up and running quickly, `martini.Classic()` provides some reasonable defaults that work well for most web applications:
~~~ go
  m := martini.Classic()
  m.Run()
~~~

Below is some of the functionality `martini.Classic()` pulls in automatically:
  * Request/Response Logging - `martini.Logger()`
  * Panic Recovery - `martini.Recovery()`
  * Static File serving - `martini.Static("public")`
  * Routing - `martini.Router`

### Handlers
Handlers are the heart and soul of Martini. A handler is basically any kind of callable function:
~~~ go
m.Get("/", func() {
  println("hello world")
}
~~~

If a handler returns a `string`, Martini will write the result to the current `*http.Request`:
~~~ go
m.Get("/", func() string {
  return "hello world" // HTTP 200 : "hello world"
})
~~~

Handlers are invoked via reflection, so Martini can inject services into Handler arguments. This makes Martini completely  compatible with golang's `http.HandlerFunc` interface:
~~~ go
m.Get("/", func(res http.ResponseWriter, req *http.Request) {
  res.WriteHead(200) // HTTP 200
})
~~~

The following services are included with a `martini.Classic()`:
  * `*log.Logger` - Global logger for Martini
  * [martini.Context](http://godoc.org/github.com/codegangsta/martini#Context) - http request context
  * `http.ResponseWriter` - http Response writer interface
  * `*http.Request` - http Request

### Routing
In Martini, a route is an HTTP method paired with a URL-matching pattern.
Each route can take one or more handler methods:
``` go
m := martini.Classic()

m.Get("/", func() {
  .. show something ..
}

m.Post("/", func() {
  .. create something ..
}

m.Put("/", func() {
  .. replace something ..
}

m.Delete("/", func() {
  .. destroy something ..
}
```

### Services
### Serving Static Files
## Middleware Handlers
### Next()
### Injecting Services

