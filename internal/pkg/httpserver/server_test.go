package httpserver

import (
	"context"
	"net/http"
	"testing"
	"time"
)

func TestRunStopsWhenContextIsCancelled(t *testing.T) {
	server := New(Config{
		Address:         "127.0.0.1:0",
		Handler:         http.NotFoundHandler(),
		ShutdownTimeout: time.Second,
	})
	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)

	go func() {
		result <- server.Run(ctx)
	}()

	cancel()

	select {
	case err := <-result:
		if err != nil {
			t.Fatalf("run server: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("server did not stop after context cancellation")
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name                  string
		config                Config
		wantReadHeaderTimeout time.Duration
		wantWriteTimeout      time.Duration
		wantIdleTimeout       time.Duration
	}{
		{
			name: "configures timeouts",
			config: Config{
				ReadHeaderTimeout: time.Second,
				WriteTimeout:      2 * time.Second,
				IdleTimeout:       3 * time.Second,
			},
			wantReadHeaderTimeout: time.Second,
			wantWriteTimeout:      2 * time.Second,
			wantIdleTimeout:       3 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := New(tt.config)

			if got := server.server.ReadHeaderTimeout; got != tt.wantReadHeaderTimeout {
				t.Fatalf("expected read header timeout %s, got %s", tt.wantReadHeaderTimeout, got)
			}
			if got := server.server.WriteTimeout; got != tt.wantWriteTimeout {
				t.Fatalf("expected write timeout %s, got %s", tt.wantWriteTimeout, got)
			}
			if got := server.server.IdleTimeout; got != tt.wantIdleTimeout {
				t.Fatalf("expected idle timeout %s, got %s", tt.wantIdleTimeout, got)
			}
		})
	}
}
