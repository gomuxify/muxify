package muxify

import (
	"errors"
	"net/http"
)

var (
	ErrNotFound = errors.New("no matching route was found")
)

type Router struct {
	NotFoundHandler http.Handler
	routes          []*Route
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var match RouteMatch
	var handler http.Handler
	if r.Match(req, &match) {
		handler = match.Handler
	}

	if handler == nil {
		handler = http.NotFoundHandler()
	}

	handler.ServeHTTP(w, req)
}

type RouteMatch struct {
	Route    *Route
	Handler  http.Handler
	MatchErr error
}

func (r *Router) NewRoute() *Route {
	route := &Route{}
	r.routes = append(r.routes, route)
	return route
}

func (r *Router) Post(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.NewRoute().Method("POST").Path(path).HandlerFunc(f)
}

func (r *Router) Get(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.NewRoute().Method("GET").Path(path).HandlerFunc(f)
}

func (r *Router) Put(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.NewRoute().Method("PUT").Path(path).HandlerFunc(f)
}

func (r *Router) Patch(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.NewRoute().Method("PATCH").Path(path).HandlerFunc(f)
}

func (r *Router) Delete(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.NewRoute().Method("DELETE").Path(path).HandlerFunc(f)
}

func (r *Router) Options(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.NewRoute().Method("OPTIONS").Path(path).HandlerFunc(f)
}

func (r *Router) Request(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.NewRoute().Path(path).HandlerFunc(f)
}

func (r *Router) Match(req *http.Request, match *RouteMatch) bool {
	for _, route := range r.routes {
		if route.Match(req, match) {
			return true
		}
	}

	if r.NotFoundHandler != nil {
		match.Handler = r.NotFoundHandler
		match.MatchErr = ErrNotFound
		return true
	}

	match.MatchErr = ErrNotFound
	return false
}
