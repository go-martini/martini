package martini

import (
	"net/http"
)

type Martini struct {
	handlers []interface{}
}

func New() *Martini {
	return &Martini{}
}

func (m *Martini) Use(handler interface{}) {
	m.handlers = append(m.handlers, handler)
}

func (m *Martini) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	context := newContext()
	context.MapTo(res, (*http.ResponseWriter)(nil))
	context.Map(req)
	context.run(m.handlers)
}
