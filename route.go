package muxify

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
)

const (
	pathParamPattern string = "([^/]+)"
)

type Route struct {
	tmplPath string
	mappings []*routeMapping
	params   []string
	matcher  *regexp.Regexp
	routeConf
}
type routeMapping struct {
	method  string
	handler http.Handler
}

func (r *Route) Match(req *http.Request, match *RouteMatch) bool {
	path := req.URL.Path

	// exact match, no need for regexp comparison
	if path == r.tmplPath {
		return r.matchHTTPMethod(req, match)
	}

	// match route by regexp comparison
	if r.matcher.MatchString(path) {

		results := r.matcher.FindAllStringSubmatch(path, 100)
		params := results[0][1:]
		for idx, param := range params {
			match.Params[r.params[idx]] = param
		}

		return r.matchHTTPMethod(req, match)
	}

	return false
}

func (r *Route) matchHTTPMethod(req *http.Request, match *RouteMatch) bool {
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

func (r *Route) Path(tmplPath string) *Route {
	r.registerMatcher(tmplPath)
	return r
}

func (r *Route) MethodHandler(method string, handler http.Handler) *Route {
	// probably a mistake to override a method handler if already set
	if r.methodExists(method) {
		err := fmt.Errorf("muxify: %s method handler for path %s already exists", method, r.tmplPath)
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

func (r *Route) registerMatcher(tmplPath string) {
	r.tmplPath = tmplPath

	s := strings.Split(r.tmplPath, "/")
	ps := []string{"^"}
	for idx, v := range s {
		if strings.HasPrefix(v, ":") {
			r.params = append(r.params, v[1:])
			ps = append(ps, pathParamPattern)
		} else {
			ps = append(ps, v)
		}

		if idx != len(s)-1 {
			ps = append(ps, "/")
		}
	}
	if r.trailingSlash {
		ps = append(ps, "/?")
	}
	ps = append(ps, "$")

	pattern := strings.Join(ps, "")
	r.matcher = regexp.MustCompile(pattern)
}
