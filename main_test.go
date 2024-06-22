package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_main(t *testing.T) {
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

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{"GET /foo", "GET", "/foo", http.StatusOK, "GET bar"},
		{"GET /foo/", "GET", "/foo/", http.StatusOK, "GET bar"},
		{"GET /foo/bar", "GET", "/foo/bar", http.StatusOK, "GET baz"},
		{"GET /foo/bar/", "GET", "/foo/bar/", http.StatusOK, "GET baz"},
		{"POST /foo/bar", "POST", "/foo/bar", http.StatusOK, "POST baz"},
		{"POST /foo/bar/", "POST", "/foo/bar/", http.StatusOK, "POST baz"},
		{"GET /invalid", "GET", "/invalid", http.StatusNotFound, "404 page not found\n"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, tt.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler := StripSlashes(mux)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.expectedStatus)
			}

			if rr.Body.String() != tt.expectedBody {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), tt.expectedBody)
			}
		})
	}
}
