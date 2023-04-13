package muxify

import (
	"fmt"
	"net/http"
)

type routeMapping struct {
	method  string
	handler http.Handler
}

type Route struct {
	path     string
	mappings []*routeMapping
}

func (r *Route) Match(req *http.Request, match *RouteMatch) bool {
	path := req.URL.Path

	if path == r.path {
		for _, mapping := range r.mappings {
			if mapping.method == req.Method {
				match.Route = r
				match.Handler = mapping.handler
				return true
			}
		}
		match.MatchErr = ErrMethodMismatch
		return false
	}

	return false
}

func (r *Route) Path(tpl string) *Route {
	r.path = tpl
	return r
}

func (r *Route) MethodHandler(method string, handler http.Handler) *Route {
	// probably a mistake to override a method handler if already set
	if r.methodExists(method) {
		err := fmt.Errorf("muxify: %s method handler for path %s already exists", method, r.path)
		panic(err)
	}
	mapping := &routeMapping{method, handler}
	r.mappings = append(r.mappings, mapping)
	return r
}

func (r *Route) MethodHandlerFunc(method string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.MethodHandler(method, http.HandlerFunc(f))
}

func (r *Route) methodExists(method string) bool {
	for _, mapping := range r.mappings {
		if mapping.method == method {
			return true
		}
	}

	return false
}
