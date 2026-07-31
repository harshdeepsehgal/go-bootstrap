// Package healthcheck contains health check handlers for the application.
package healthcheck

import (
	"encoding/json"
	"net/http"
)

const statusOK = "ok"

// Report is the JSON response from the health endpoint.
type Report struct {
	Status string `json:"status"`
}

// HealthCheckHandler reports that the application process is healthy.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(Report{Status: statusOK})
}
