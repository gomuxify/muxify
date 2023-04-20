package muxify

import (
	"context"
	"errors"
	"net/http"
)

var (
	ErrNotFound       = errors.New("no matching route was found")
	ErrMethodMismatch = errors.New("method is not allowed")
)

type Router struct {
	NotFoundHandler         http.Handler
	MethodNotAllowedHandler http.Handler
	routes                  []*Route
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var match RouteMatch
	var handler http.Handler
	if r.Match(req, &match) {
		handler = match.Handler
		req = ctxAddParams(req, match.Params)
	}

	if handler == nil && match.MatchErr == ErrMethodMismatch {
		handler = methodNotAllowedHandler()
	}

	if handler == nil {
		handler = http.NotFoundHandler()
	}

	handler.ServeHTTP(w, req)
}

type RouteMatch struct {
	Route    *Route
	Handler  http.Handler
	Params   map[string]string
	MatchErr error
}

func (r *Router) NewRoute() *Route {
	route := &Route{paramPos: make(map[string]int)}
	r.routes = append(r.routes, route)
	return route
}

func (r *Router) Post(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.upsertPath(path, "POST", f)
}

func (r *Router) Get(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.upsertPath(path, "GET", f)
}

func (r *Router) Put(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.upsertPath(path, "PUT", f)
}

func (r *Router) Patch(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.upsertPath(path, "PATCH", f)
}

func (r *Router) Delete(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.upsertPath(path, "DELETE", f)
}

func (r *Router) Options(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.upsertPath(path, "OPTIONS", f)
}

func (r *Router) Match(req *http.Request, match *RouteMatch) bool {
	for _, route := range r.routes {
		if route.Match(req, match) {
			return true
		}
	}

	if match.MatchErr == ErrMethodMismatch {
		if r.MethodNotAllowedHandler != nil {
			match.Handler = r.MethodNotAllowedHandler
			return true
		}

		return false
	}

	if r.NotFoundHandler != nil {
		match.Handler = r.NotFoundHandler
		match.MatchErr = ErrNotFound
		return true
	}

	match.MatchErr = ErrNotFound
	return false
}

func Params(r *http.Request) map[string]string {
	if rv := r.Context().Value(paramsKey); rv != nil {
		return rv.(map[string]string)
	}
	return nil
}

func (r *Router) upsertPath(path, method string, f func(http.ResponseWriter, *http.Request)) *Route {
	if route := r.findRouteWithPath(path); route != nil {
		return route.MethodHandlerFunc(method, f)
	}

	return r.NewRoute().Path(path).MethodHandlerFunc(method, f)
}

func (r *Router) findRouteWithPath(path string) *Route {
	for _, route := range r.routes {
		if route.tmplPath == path {
			return route
		}
	}

	return nil
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func methodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(methodNotAllowed)
}

type contextKey int

const (
	paramsKey contextKey = iota
)

func ctxAddParams(r *http.Request, params map[string]string) *http.Request {
	ctx := context.WithValue(r.Context(), paramsKey, params)
	return r.WithContext(ctx)
}
