package recovery

import (
  "github.com/codegangsta/martini"
)

func New() martini.Handler {
  return func(c martini.Context) {
    defer handlePanic()
    c.Next()
  }
}

func handlePanic() {
  if err := recover(); err != nil {
    println(err)
  }
}
