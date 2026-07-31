// Package middleware provides small, dependency-free HTTP middleware helpers.
package middleware

import (
	"net/http"
	"runtime/debug"

	"go-bootstrap/internal/pkg/logger"
	"go.uber.org/zap"
)

// Middleware wraps an HTTP handler.
type Middleware func(http.Handler) http.Handler

// Chain applies middleware in declaration order.
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
	for index := len(middlewares) - 1; index >= 0; index-- {
		handler = middlewares[index](handler)
	}
	return handler
}

// Recover logs panics and returns an internal-server-error response.
func Recover() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					logger.Logger().Error(
						"panic while handling request",
						zap.Any("panic", recovered),
						zap.String("method", r.Method),
						zap.String("path", r.URL.Path),
						zap.ByteString("stack", debug.Stack()),
					)
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
