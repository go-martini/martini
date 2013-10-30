package martini

import (
	"github.com/codegangsta/inject"
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
	ctx := &context{inject.New(), m.handlers, 0}
	ctx.MapTo(ctx, (*Context)(nil))
	ctx.MapTo(res, (*http.ResponseWriter)(nil))
	ctx.Map(req)
	ctx.run()
}
