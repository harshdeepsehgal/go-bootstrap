package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go-bootstrap/internal/pkg/logger"

	"go.uber.org/zap"
)

func TestRecover(t *testing.T) {
	setTestLogger(t, zap.NewNop())

	panicMiddleware := func(http.Handler) http.Handler {
		return http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
			panic("boom")
		})
	}
	tests := []struct {
		name        string
		handler     http.Handler
		middlewares []Middleware
		wantStatus  int
	}{
		{
			name: "recovers handler panic",
			handler: http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
				panic("boom")
			}),
			middlewares: []Middleware{Recover()},
			wantStatus:  http.StatusInternalServerError,
		},
		{
			name:        "recovers middleware panic",
			handler:     http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}),
			middlewares: []Middleware{Recover(), panicMiddleware},
			wantStatus:  http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := httptest.NewRecorder()
			Chain(tt.handler, tt.middlewares...).ServeHTTP(
				response,
				httptest.NewRequest(http.MethodGet, "/", nil),
			)

			if response.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, response.Code)
			}
		})
	}
}

func TestContext(t *testing.T) {
	tests := []struct {
		name          string
		requestID     string
		wantRequestID string
	}{
		{
			name:          "preserves supplied request ID",
			requestID:     "request-123",
			wantRequestID: "request-123",
		},
		{
			name: "generates missing request ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var contextRequestID string
			handler := Context()(http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				contextRequestID, _ = logger.RequestID(r.Context())
			}))
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.requestID != "" {
				request.Header.Set("X-Request-ID", tt.requestID)
			}
			response := httptest.NewRecorder()

			handler.ServeHTTP(response, request)

			responseRequestID := response.Header().Get("X-Request-ID")
			if responseRequestID == "" {
				t.Fatal("expected response request ID")
			}
			if responseRequestID != contextRequestID {
				t.Fatalf("expected matching response and context request IDs, got %q and %q", responseRequestID, contextRequestID)
			}
			if tt.wantRequestID != "" && responseRequestID != tt.wantRequestID {
				t.Fatalf("expected request ID %q, got %q", tt.wantRequestID, responseRequestID)
			}
		})
	}
}

func setTestLogger(t *testing.T, testLogger *zap.Logger) {
	t.Helper()
	originalLogger := logger.Log
	logger.Log = testLogger
	t.Cleanup(func() {
		logger.Log = originalLogger
	})
}
