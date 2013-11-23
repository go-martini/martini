package martini

import (
    "strconv"
)

type RouteHelper struct {
    Router Router
}

// UrlFor returns the url for the given route name.
func (rh *RouteHelper) UrlFor(routeName string, params ...interface{}) string {
	var args []string
	for _, param := range params {
		switch v := param.(type) {
		case int:
			args = append(args, strconv.FormatInt(int64(v), 10))
		case string:
			args = append(args, v)
		default:
			if v != nil {
				panic("Arguments passed to UrlFor must be integers or strings")
			}
		}
	}

	for _, route := range rh.Router.GetRoutes() {
		if route.GetName() == routeName {
			return route.UrlWith(args)
		}
	}

	return ""
}