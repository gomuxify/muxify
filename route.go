package muxify

import (
	"net/http"
)

type Route struct {
	path    string
	method  string
	handler http.Handler
}

func (r *Route) Match(req *http.Request, match *RouteMatch) bool {
	path := req.URL.Path

	if path == r.path && (r.method == "" || req.Method == r.method) {
		match.Route = r
		match.Handler = r.handler
		return true
	}

	return false
}

func (r *Route) Path(tpl string) *Route {
	r.path = tpl
	return r
}

func (r *Route) Method(method string) *Route {
	r.method = method
	return r
}

func (r *Route) Handler(handler http.Handler) *Route {
	r.handler = handler
	return r
}

func (r *Route) HandlerFunc(f func(http.ResponseWriter, *http.Request)) *Route {
	return r.Handler(http.HandlerFunc(f))
}
