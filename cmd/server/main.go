package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"go-bootstrap/config"
	"go-bootstrap/internal/go-bootstrap/handler"
	"go-bootstrap/internal/pkg/healthcheck"
	"go-bootstrap/internal/pkg/httpserver"
	"go-bootstrap/internal/pkg/logger"
	"go-bootstrap/internal/pkg/middleware"

	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "go-bootstrap:", err)
		os.Exit(1)
	}
}

func run() error {
	config.Initialize()
	cfg := config.GetConfig()

	logger.Init()

	httpHandler := NewHTTPHandler()
	httpServer := httpserver.New(httpserver.Config{
		Address:           cfg.HTTP.Address,
		Handler:           httpHandler,
		ReadHeaderTimeout: cfg.HTTP.ReadHeaderTimeout,
		WriteTimeout:      cfg.HTTP.WriteTimeout,
		IdleTimeout:       cfg.HTTP.IdleTimeout,
		ShutdownTimeout:   cfg.HTTP.ShutdownTimeout,
	})
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := httpServer.Run(ctx); err != nil {
		logger.Logger().Error("HTTP server stopped unexpectedly", zap.Error(err))
		return err
	}

	return nil
}

// NewHTTPHandler builds the application's HTTP routes and middleware.
func NewHTTPHandler() http.Handler {
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)
	mux.HandleFunc("GET /healthz", healthcheck.HealthCheckHandler)

	return middleware.Chain(mux, middleware.Recover())
}
