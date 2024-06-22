# muxtf

I used only `net/http` in a recent project and found this unexpected behavior:

- `POST /endpoint/` would invoke the expected HandlerFunc`
- `POST /endoint` would fall through to `GET /endpoint`

The goal was to get the proper response no matter if the URL path has a trailing slash.

The route group created by `NewRouterGroup` in combination with the `StripSlashes` middleware is what I came up with to achieve
what I wanted. I am unsure if this is the idiomatic Go way of doing it... Anyways, the test are passing ¯\_(ツ)_/¯

Any feedback would is highly appreciated. This thing feels like a massive skill issue and I want to get rid of it!

# Initial Problem

This is how I created "route groups":

```go
employeeHandlers := newEmployeeHandler(ah)
mux.Handle("/employees/", http.StripPrefix("/employees", employeeHandlers.router))
```

```go
func newEmployeeHandler(ah *APIHandler) *employeeHandler {
    h := &employeeHandler{
        ah:      ah,
        service: service.NewEmployeeService(ah.db, ah.ug),
        router:  http.NewServeMux(),
    }

    h.router.HandleFunc("PUT /link/{id}", h.linkNfc)

    // Add protected handlers here
    protectedMux := http.NewServeMux()
    protectedMux.HandleFunc("DELETE /{id}", h.deleteEmployee)
    protectedMux.HandleFunc("GET /link/{id}", h.requestQRPayload)
    protectedMux.HandleFunc("POST /", h.createEmployee)
    protectedMux.HandleFunc("GET /", h.getEmployeeList)

    h.router.Handle("/", middleware.NewAuthenticationMiddleware(h.ah.authService)(protectedMux))

    return h
}
```

Doing so, would exhibit the unwanted behavior described above.
