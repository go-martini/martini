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
`go get github.com/codegangsta/martini`

View at: http://localhost:3000

## Table of contents
* [Martini](#martini)

