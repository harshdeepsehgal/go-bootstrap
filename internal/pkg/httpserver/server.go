// Package httpserver provides a small, graceful HTTP server wrapper.
package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// Config configures an HTTP server.
type Config struct {
	Address           string
	Handler           http.Handler
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

// Server runs an HTTP server and shuts it down when its context is cancelled.
type Server struct {
	shutdownTimeout time.Duration
	server          *http.Server
}

// New creates a server from config.
func New(config Config) *Server {
	return &Server{
		shutdownTimeout: config.ShutdownTimeout,
		server: &http.Server{
			Addr:              config.Address,
			Handler:           config.Handler,
			ReadHeaderTimeout: config.ReadHeaderTimeout,
			WriteTimeout:      config.WriteTimeout,
			IdleTimeout:       config.IdleTimeout,
		},
	}
}

// Run listens on the configured address and serves requests until the server
// fails or ctx is cancelled.
func (s *Server) Run(ctx context.Context) error {
	errs := make(chan error, 1)
	go func() {
		errs <- s.server.ListenAndServe()
	}()

	select {
	case err := <-errs:
		return serveError(err)
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
		defer cancel()

		if err := s.server.Shutdown(shutdownCtx); err != nil {
			_ = s.server.Close()
			<-errs
			return fmt.Errorf("shut down HTTP server: %w", err)
		}

		return serveError(<-errs)
	}
}

func serveError(err error) error {
	if err == nil || errors.Is(err, http.ErrServerClosed) {
		return nil
	}

	return fmt.Errorf("serve HTTP server: %w", err)
}
