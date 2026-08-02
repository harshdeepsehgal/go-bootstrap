// Package logger configures and exposes the process-wide structured logger.
package logger

import (
	"context"
	"fmt"
	"go-bootstrap/config"
	"log"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logOnce sync.Once
	Log     *zap.Logger
)

type contextKey uint8

const requestIDKey contextKey = iota

// Init configures the process-wide logger. The first call wins and must happen
// during process startup before serving requests.
func Init() {
	logOnce.Do(func() {
		configuredLogger, err := newLogger(config.GetConfig().Log.Level)
		if err != nil {
			log.Fatal("failed to initialise logger")
		}

		Log = configuredLogger
	})
}

func newLogger(level string) (*zap.Logger, error) {
	level = strings.TrimSpace(level)
	if level == "" {
		level = "info"
	}

	zapConfig := zap.NewProductionConfig()
	zapConfig.Level = parseLevel(level)
	zapConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	configuredLogger, err := zapConfig.Build()
	if err != nil {
		return nil, fmt.Errorf("build logger: %w", err)
	}

	return configuredLogger, nil
}

func parseLevel(level string) zap.AtomicLevel {
	switch level {
	case "info":
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	case "error":
		return zap.NewAtomicLevelAt(zap.ErrorLevel)
	case "debug":
		return zap.NewAtomicLevelAt(zap.DebugLevel)
	case "warn":
		return zap.NewAtomicLevelAt(zap.WarnLevel)
	default:
		return zap.NewAtomicLevelAt(zap.InfoLevel)
	}
}

// WithRequestID stores the request ID in the context.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey, requestID)
}

// RequestID extracts the request ID from the context.
func RequestID(ctx context.Context) (string, bool) {
	if ctx == nil {
		return "", false
	}

	requestID, ok := ctx.Value(requestIDKey).(string)
	return requestID, ok && requestID != ""
}

// Logger returns the initialized process-wide logger.
func Logger(ctx context.Context) *zap.Logger {
	if Log == nil {
		return zap.NewNop()
	}
	if requestID, ok := RequestID(ctx); ok {
		return Log.With(zap.String("rqId", requestID))
	}

	return Log
}
