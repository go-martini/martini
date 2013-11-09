# Martini [![Build Status](https://drone.io/github.com/codegangsta/martini/status.png)](https://drone.io/github.com/codegangsta/martini/latest)

Martini is a package for quickly writing modular and maintainable web applications/services in Golang.

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

View at: http://localhost:3000

## Table of Contents
* [Martini](#martini-)
  * [Table of Contents](#table-of-contents)
  * [Classic Martini](#classic-martini)
    * [Routing](#routing)
    * [Handlers](#handlers)
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

### Routing
### Handlers
### Services
### Serving Static Files
## Middleware Handlers
### Next()
### Injecting Services

