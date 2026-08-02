package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRegisterRoutes(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
	}{
		{
			name:       "allows example users route",
			method:     http.MethodGet,
			path:       "/users",
			wantStatus: http.StatusOK,
		},
		{
			name:       "rejects unsupported method",
			method:     http.MethodPost,
			path:       "/users",
			wantStatus: http.StatusMethodNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			testHandler().ServeHTTP(
				response,
				httptest.NewRequest(tt.method, tt.path, nil),
			)

			if response.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, response.Code)
			}
		})
	}
}

func testHandler() http.Handler {
	mux := http.NewServeMux()
	RegisterRoutes(mux)
	return mux
}
