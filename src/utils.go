package goway

import (
	"context"
	"net/http"
	"path"
	"strings"
)

// when ServeHTTP called, sets Executed to true.
type dummyEndpoint struct {
	Executed bool
}

func (d *dummyEndpoint) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	d.Executed = true
}

// run middlewares.
func executeMiddleware(response http.ResponseWriter, request *http.Request, middleware MiddlewareFunc) (isResponseSended bool) {
	if middleware == nil {
		return
	}

	// set dummy as middleware endpoint.
	var dummy = &dummyEndpoint{}
	var handler = middleware(dummy)
	handler.ServeHTTP(response, request)

	// if dummy not executed
	// it means middleware sended response / not called next.ServeHTTP()
	isResponseSended = !dummy.Executed

	return
}

// https://gist.github.com/husobee/fd23681261a39699ee37?permalink_comment_id=3111569#gistcomment-3111569
//
// make one big middleware from middlewares.
func middlewareChain(middlewares ...MiddlewareFunc) MiddlewareFunc {
	if len(middlewares) < 1 {
		return nil
	}
	return func(next http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

// get path variables.
func Vars(request *http.Request) map[string]string {
	var ctx = request.Context()
	var varsMap, ok = ctx.Value(CTX_VARS_NAME).(map[string]string)
	if !ok {
		return nil
	}
	return varsMap
}

// is route variable?
func isRouteVar(path string) (isVar bool, varName string) {
	isVar = strings.HasPrefix(path, "{") && strings.HasSuffix(path, "}")
	if isVar {
		var withoutBrackets = strings.ReplaceAll(path, "{", "")
		withoutBrackets = strings.ReplaceAll(withoutBrackets, "}", "")
		varName = withoutBrackets
	}
	return
}

// add route variable with value to request context.
func addVarToContext(request *http.Request, name string, value string) {
	var oldCtx = request.Context()
	var varsMap, ok = oldCtx.Value(CTX_VARS_NAME).(map[string]string)
	if !ok || varsMap == nil {
		varsMap = make(map[string]string, 0)
	}
	varsMap[name] = value
	var ctxWithVars = context.WithValue(oldCtx, CTX_VARS_NAME, varsMap)
	*request = *request.WithContext(ctxWithVars)
}

// check is method allowed.
func isMethodAllowed(methods []string, requestMethod string) bool {
	if methods == nil {
		return true
	}
	for _, method := range methods {
		if requestMethod == method {
			return true
		}
	}
	return false
}

// remove slash at start and end of str.
func removeSlashStartEnd(str string) string {
	if len(str) < 1 {
		return str
	}
	str = removeSlashStart(str)
	str = removeSlashEnd(str)
	str = path.Clean(str)
	return str
}

// remove slash at start of str.
func removeSlashStart(str string) string {
	str = strings.TrimPrefix(str, "/")
	return str
}

// remove slash at end of str.
func removeSlashEnd(str string) string {
	str = strings.TrimSuffix(str, "/")
	return str
}

// default 405 handler.
func getDefaultHandler405() RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(405)
		w.Write([]byte("method not allowed"))
	}
}

// default 404 handler.
func getDefaultHandler404() RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("not found"))
	}
}

// make path like: /hello/world
func pathToStandart(to string) string {
	if len(to) < 1 {
		return to
	}
	to = "/" + removeSlashEnd(to)
	to = path.Clean(to)
	if to == "." || to == "/" {
		to = ""
	}
	return to
}

// https://www.geeksforgeeks.org/how-to-remove-duplicate-values-from-slice-in-golang/
//
// remove duplicates from slice.
func removeDuplicateValues(slice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	// If the key(values of the slice) is not equal
	// to the already present value in new slice (list)
	// then we append it. else we jump on another element.
	for _, entry := range slice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

// format route allowed methods.
func processAllowedMethods(slice []string, methods ...string) []string {
	if slice == nil {
		slice = make([]string, 0)
	}
	for _, method := range methods {
		method = strings.ToUpper(method)
		slice = append(slice, method)
	}
	slice = removeDuplicateValues(slice)
	return slice
}

func processMiddleware(old MiddlewareFunc, middlewares ...MiddlewareFunc) MiddlewareFunc {
	var chained = make([]MiddlewareFunc, 0)
	if old != nil {
		chained = append(chained, old)
	}
	for _, m := range middlewares {
		if m == nil {
			continue
		}
		chained = append(chained, m)
	}
	var finalChain = middlewareChain(chained...)
	return finalChain
}

func splitPath(path string) []string {
	path = pathToStandart(path)
	if len(path) < 1 {
		return make([]string, 0)
	}
	return strings.Split(removeSlashStartEnd(path), "/")
}

func isPathSliceEmpty(slice []string) bool {
	return slice == nil || len(slice) < 1 || slice[0] == "" || slice[0] == "."
}
