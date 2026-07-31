// Package handler owns the application's HTTP handlers and routes.
package handler

import (
	"net/http"
)

// RegisterRoutes registers the application's HTTP routes, including an example users endpoint.
func RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /users", exampleGetUsersHandler)
}

// exampleGetUsersHandler demonstrates a minimal application endpoint.
func exampleGetUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("List of users"))
}
