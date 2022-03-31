package goway

import (
	"net/http"
)

type Route struct {
	// is route under group?
	isUnderGroup bool

	// prefix tools.
	prefix prefixes

	// allowed route methods.
	allowedMethods []string

	// route middleware chain.
	middleware MiddlewareFunc

	// route endpoint.
	handler RouteHandler
}

func (r *Route) new(excludeCount int, to string, handler RouteHandler) {
	if excludeCount > 0 {
		r.isUnderGroup = true
		r.prefix.excludeCount = excludeCount
	}
	r.prefix.setPath(to)
	r.prefix.setPathSlice()
	r.handler = handler
}

func (r *Route) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// run middleware.
	var isResponseSended = executeMiddleware(response, request, r.middleware)
	if isResponseSended {
		return
	}

	r.handler(response, request)
	return
}

// route trigger on this methods only.
func (r *Route) Methods(methods ...string) *Route {
	r.allowedMethods = processAllowedMethods(r.allowedMethods, methods...)
	return r
}

// provide middleware.
func (r *Route) Use(middleware ...MiddlewareFunc) *Route {
	r.middleware = processMiddleware(r.middleware, middleware...)
	return r
}
