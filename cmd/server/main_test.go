package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHTTPHandlerAllowsHealthCheck(t *testing.T) {
	response := httptest.NewRecorder()
	NewHTTPHandler().ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/healthz", nil),
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestHTTPHandlerRejectsUnsupportedHealthCheckMethod(t *testing.T) {
	response := httptest.NewRecorder()
	NewHTTPHandler().ServeHTTP(
		response,
		httptest.NewRequest(http.MethodPost, "/healthz", nil),
	)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", response.Code)
	}
}
