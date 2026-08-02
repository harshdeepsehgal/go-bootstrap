package logger

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name               string
		ctx                context.Context
		wantRequestID      string
		wantRequestIDField bool
	}{
		{
			name:               "adds request ID from context",
			ctx:                WithRequestID(context.Background(), "request-123"),
			wantRequestID:      "request-123",
			wantRequestIDField: true,
		},
		{
			name: "does not add empty request ID",
			ctx:  context.Background(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			core, logs := observer.New(zap.InfoLevel)
			originalLogger := Log
			Log = zap.New(core)
			t.Cleanup(func() {
				Log = originalLogger
			})

			Logger(tt.ctx).Info("test message")

			entries := logs.All()
			if len(entries) != 1 {
				t.Fatalf("expected one log entry, got %d", len(entries))
			}
			requestID, exists := entries[0].ContextMap()["rqId"]
			if exists != tt.wantRequestIDField {
				t.Fatalf("expected request ID field presence %t, got %t", tt.wantRequestIDField, exists)
			}
			if exists && requestID != tt.wantRequestID {
				t.Fatalf("expected request ID %s, got %v", tt.wantRequestID, requestID)
			}
		})
	}
}
