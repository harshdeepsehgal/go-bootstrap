package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExampleUsersRouteAllowsGET(t *testing.T) {
	response := httptest.NewRecorder()
	testHandler().ServeHTTP(
		response,
		httptest.NewRequest(http.MethodGet, "/users", nil),
	)

	if response.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", response.Code)
	}
}

func TestExampleUsersRouteRejectsUnsupportedMethod(t *testing.T) {
	response := httptest.NewRecorder()
	testHandler().ServeHTTP(
		response,
		httptest.NewRequest(http.MethodPost, "/users", nil),
	)

	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected status 405, got %d", response.Code)
	}
}

func testHandler() http.Handler {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	return mux
}
