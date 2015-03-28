package main

import (
	"github.com/go-martini/martini"
)

func main() {
	m := martini.Classic()
	m.Get("/", func() string {
		return "Hello world test file!"
	})

	m.Get("/hi", func() string {
		return "Say Hi!"
	})
	// m.NotFound(func() string {
	// 	return "test"
	// })
	m.Run()
}
