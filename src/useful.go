package goway

import (
	"net/http"
)

// middleware.
type MiddlewareFunc func(http.Handler) http.Handler

// route endpoint.
type RouteHandler func(http.ResponseWriter, *http.Request)

// vars for request context.
type CTX_VAL string

const (
	// route variables.
	CTX_VARS_NAME CTX_VAL = "GOWAY_ROUTER_VARS"
)

// when route not found.
var Handler404 = getDefaultHandler404()

// when request method not allowed.
var Handler405 = getDefaultHandler405()

// tools for working on route/group paths.
type prefixes struct {
	// because we dealing with nested routing
	// paths can be computed before (in route groups for example).
	// This count displays how much elements we need to cut from request path
	// for route matching.
	//
	// example:
	//
	// route group path slice: [api, users]
	//
	// request path slice: [api, users, me]
	//
	// group computed before, and excludeCount equals 2
	//
	// cut first two elements from request path slice
	// and we get: [me] - now we can match route
	excludeCount int

	// path or group prefix like: /hello/world.
	path string

	// path slice like: [hello, world].
	pathSlice []string
}

// exclude path by requestPath and prefix exclude.
func (p *prefixes) getExcluded(requestPath string) []string {
	if len(requestPath) < 1 {
		return make([]string, 0)
	}
	var splitted = splitPath(requestPath)
	if len(splitted) >= p.excludeCount {
		return splitted[p.excludeCount:]

	}
	return splitted
}

// set path.
func (p *prefixes) setPath(to string) {
	p.path = pathToStandart(to)
}

// set pathSlice by path.
func (p *prefixes) setPathSlice() {
	var splitted = splitPath(p.path)
	p.pathSlice = splitted
}
