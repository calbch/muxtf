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
			handler: foo,
		},
		{
			path:    "GET /baz",
			handler: bazGET,
		},

		{
			path:    "POST /baz",
			handler: bazPOST,
		},
	}

	NewRouteGroup("test", GroupParams{
		handlerList: handlers,
	}, mux)

	err := http.ListenAndServe(":8080", StripSlashes(mux))
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}

func foo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("foo"))
}

func bazGET(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GET request to %s\n", r.URL)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("GET"))
}

func bazPOST(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("POST request to %s\n", r.URL)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("POST"))
}

type Handler struct {
	path    string
	handler http.HandlerFunc
}

type GroupParams struct {
	handlerList []Handler
}

func NewRouteGroup(prefix string, arg GroupParams, r *http.ServeMux) {
	fmt.Printf("Creating Route Group '%s'\n", strings.ToUpper(prefix))
	for _, h := range arg.handlerList {
		path := addPrefix(prefix, h.path)

		r.HandleFunc(path, h.handler)
		fmt.Println(path)
	}
	fmt.Println("-----------DONE")
}

func addPrefix(prefix, p string) string {
	splits := strings.Split(p, " ")

	// add validation here

	if len(splits) > 1 {
		return fmt.Sprintf("%s /%s%s", splits[0], prefix, splits[1])
	}
	return fmt.Sprintf("%s /%s", splits[0], prefix)

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
