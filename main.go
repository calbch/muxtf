package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()

	handlers := []Handler{
		{
			path:    "GET",
			handler: getBarHandler,
		},
		{
			path:    "GET /bar",
			handler: getBazHandler,
		},
		{
			path:    "POST /bar",
			handler: postBazHandler,
		},
	}
	NewRouteGroup("foo", handlers, mux)

	handler := StripSlashes(mux)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}

// THE IMPORTANT PART IS IN HERE
type Handler struct {
	path    string
	handler http.HandlerFunc
}

// Creates a prefixed route group on a passed http.ServeMux
func NewRouteGroup(prefix string, handlerList []Handler, mux *http.ServeMux) {
	for _, h := range handlerList {
		path := addPrefix(prefix, h.path)
		mux.HandleFunc(path, h.handler)
	}
}

// Build the URL path depening on wether a subpath is used or not
func addPrefix(prefix, p string) string {
	parts := strings.SplitN(p, " ", 2)
	method := parts[0]
	if len(parts) > 1 {
		path := parts[1]
		return fmt.Sprintf("%s /%s%s", method, prefix, path)
	}
	return fmt.Sprintf("%s /%s", method, prefix)
}

func StripSlashes(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		l := len(p)
		if l > 1 && p[l-1] == '/' {
			r.URL.Path = p[:l-1]
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// END

func getBarHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("GET bar"))
}

func getBazHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("GET baz"))
}

func postBazHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("POST baz"))
}
