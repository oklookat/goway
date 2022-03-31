# goway — golang router


## Features
- Route groups
- Allowed methods
- Middlewares
- Custom 404/405 handler


## Example


```go 
import (
    "net/http"
    "github.com/oklookat/goway"
)

var rootMiddleware = func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		next.ServeHTTP(response, request)
		return
	})
}

func main() {
    var root = goway.New()
    root.Use(rootMiddleware)

    var group = root.Group("/api")
    group.Route("/users", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "/api/users route")
	}).Methods(http.MethodGet)

    root.Route("/another", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "/another route")
	})

    http.Handle("/", router)
}
```

