package goway

import (
	"net/http"
)

type routeMatcher struct {
	request          *http.Request
	requestPath      string
	requestPathSlice []string
	method           string
}

func (r *routeMatcher) New(req *http.Request) {
	r.request = req
	r.method = r.request.Method

	// convert request path to standart path, like we do it with route paths.
	r.requestPath = pathToStandart(req.URL.Path)
	r.requestPathSlice = splitPath(r.requestPath)
}

// match route group. Returns:
//
// group / statusCode 0
//
// nil / statusCode 404/405
func (r *routeMatcher) Groups(routers []*Router) (matched *Router, statusCode int) {
	statusCode = 200
	var requestPath = r.requestPath
	var requestPathSlice = r.requestPathSlice

	for i := range routers {
		// check is route has parent prefix.
		var isHasExclude = routers[i].prefix.excludeCount > 0
		if isHasExclude {
			requestPathSlice = routers[i].prefix.getExcluded(requestPath)
		}

		// example: prefix /hello/world, request /hello. Not our group.
		if len(routers[i].prefix.pathSlice) > len(r.requestPathSlice) {
			continue
		}

		var isPiecesMatched = r.matchPathPieces(i, routers[i].prefix.pathSlice, requestPathSlice)
		if isPiecesMatched {
			matched = routers[i]
		}

		var isLastGroup = i == len(routers)-1

		// no match.
		if matched == nil {
			if !isLastGroup {
				// not last group, continue searching.
				continue
			}
			// change code to 404, but if before we not found match with unallowed method.
			if statusCode != 405 {
				statusCode = 404
			}
			return
		}

		// match.
		// check is method allowed.
		var isAllowed = isMethodAllowed(matched.allowedMethods, r.method)
		if isAllowed {
			// it's our match.
			statusCode = 0
			return
		}
		// method not allowed.
		statusCode = 405
		if isLastGroup {
			// method not allowed and no more groups.
			matched = nil
			return
		}

		// summary: group matched, but method not allowed
		// group not last, try to find other group.
	}
	return
}

// match route. Returns:
//
// route / statusCode 0
//
// nil / statusCode 404/405
func (r *routeMatcher) Routes(routes []*Route) (matched *Route, statusCode int) {
	statusCode = 200
	var requestPath = r.requestPath
	var requestPathSlice = r.requestPathSlice
	for i := range routes {
		// check is route under group (has prefix).
		if routes[i].isUnderGroup {
			// remove prefix from request path (its already computed in group).
			requestPathSlice = routes[i].prefix.getExcluded(requestPath)
		}

		// example: request /hello, route /hello/world. Not our route.
		if len(routes[i].prefix.pathSlice) != len(requestPathSlice) {
			continue
		}

		var isPiecesMatched = r.matchPathPieces(i, routes[i].prefix.pathSlice, requestPathSlice)
		if isPiecesMatched {
			matched = routes[i]
		}

		var isLastRoute = i == len(routes)-1

		// no match.
		if matched == nil {
			if !isLastRoute {
				// not last route, continue searching.
				continue
			}
			// change code to 404, but if before we not found match with unallowed method.
			if statusCode != 405 {
				statusCode = 404
			}
			return
		}

		// match.
		// check is method allowed.
		var isAllowed = isMethodAllowed(matched.allowedMethods, r.method)
		if isAllowed {
			// it's our match.
			statusCode = 0
			return
		}
		// method not allowed.
		statusCode = 405
		if isLastRoute {
			// method not allowed and no more routes.
			matched = nil
			return
		}

		// summary: route matched, but method not allowed
		// route not last, try to find other route.

	}
	return
}

// compare paths and add route vars to request context if exists.
//
// Matched == true examples:
//
// 1. pathSlice: ["api", "users", "{id}"], requestPathSlice ["api", "users", "12", "actions", "hello"]
//
// 2. pathSlice: ["api", "users", "12"], requestPathSlice ["api", "users", "12", "actions"]
//
// 3. pathSlice: ["api", "users"], requestPathSlice ["api", "users"]
//
// 4. pathSlice: [] or nil, requestPathSlice [] or nil
func (r *routeMatcher) matchPathPieces(counter int, pathSlice []string, requestPathSlice []string) (matched bool) {
	// match if paths empty.
	if isPathSliceEmpty(pathSlice) && isPathSliceEmpty(requestPathSlice) {
		return true
	}

	// compare paths.
	for pieceCounter := range pathSlice {
		var pathPiece = pathSlice[pieceCounter]
		var requestPathPiece = requestPathSlice[pieceCounter]

		var isSame = pathPiece == requestPathPiece
		// not same? maybe it's route variable?
		if !isSame {
			var isVar, name = isRouteVar(pathPiece)
			if !isVar {
				// summary: pieces not same, and it's not var. Not our path.
				matched = false
				break
			}

			// if it's var, add it to context
			addVarToContext(r.request, name, requestPathPiece)
		}

		// check is last piece.
		var isLast = pieceCounter == len(pathSlice)-1
		if isLast {
			// pieces same, maybe it's our path.
			matched = true
			return
		}
	}
	return
}
