package goway

import (
	"net/http"
)

/*
goway router
some ideas from gorilla/mux router
https://github.com/gorilla/mux/blob/master/LICENSE
*/

// create new router.
func New() *Router {
	return &Router{}
}

type Router struct {
	// route groups.
	groups []*Router

	// routes.
	routes []*Route

	// prefix info.
	prefix prefixes

	// allowed request methods.
	allowedMethods []string

	// middleware chain.
	middleware MiddlewareFunc
}

// any parents (routes or groups) should remove this exclude prefix from
// request path to match request.
func (r *Router) getExcludePrefix() (excludeCount int) {
	if len(r.prefix.path) < 1 {
		return 0
	}
	var exclude = "/" + r.prefix.path
	exclude = pathToStandart(exclude)
	var excludeSlice = splitPath(exclude)
	excludeCount = len(excludeSlice) + r.prefix.excludeCount
	return
}

// when request coming.
func (r *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	// run middleware.
	var isResponseSended = executeMiddleware(response, request, r.middleware)
	if isResponseSended {
		return
	}

	var matcher = routeMatcher{}
	matcher.New(request)

	// try to match groups first.
	if r.groups != nil {
		var matched, code = matcher.Groups(r.groups)
		if code == 405 {
			Handler405(response, request)
			return
		}
		if matched != nil {
			matched.ServeHTTP(response, request)
			return
		}
	}

	// try to match routes.
	if r.routes != nil {
		var matched, code = matcher.Routes(r.routes)
		if code == 405 {
			Handler405(response, request)
			return
		}
		if matched != nil {
			matched.ServeHTTP(response, request)
			return
		}
	}

	// 404.
	Handler404(response, request)
	return
}

// add route.
func (r *Router) Route(to string, handler RouteHandler) *Route {
	if r.routes == nil {
		r.routes = make([]*Route, 0)
	}

	// new route.
	var newRoute = &Route{}
	var excludeCount = r.getExcludePrefix()
	newRoute.new(excludeCount, to, handler)

	// add to routes.
	r.routes = append(r.routes, newRoute)
	return newRoute
}

// add route group.
func (r *Router) Group(prefix string) (group *Router) {
	// make groups if not.
	if r.groups == nil {
		r.groups = make([]*Router, 0)
	}

	// create new router.
	var newRouter = New()
	var excludeCount = r.getExcludePrefix()
	newRouter.prefix.excludeCount = excludeCount
	newRouter.prefix.setPath(prefix)
	newRouter.prefix.setPathSlice()

	// add to groups.
	r.groups = append(r.groups, newRouter)
	return newRouter
}

// provide middleware.
func (r *Router) Use(middleware ...MiddlewareFunc) *Router {
	r.middleware = processMiddleware(r.middleware, middleware...)
	return r
}

// add allowed request methods.
func (r *Router) Methods(methods ...string) *Router {
	r.allowedMethods = processAllowedMethods(r.allowedMethods, methods...)
	return r
}
