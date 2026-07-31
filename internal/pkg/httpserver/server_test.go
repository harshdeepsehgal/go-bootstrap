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

func TestNewConfiguresTimeouts(t *testing.T) {
	server := New(Config{
		ReadHeaderTimeout: time.Second,
		WriteTimeout:      2 * time.Second,
		IdleTimeout:       3 * time.Second,
	})

	if server.server.ReadHeaderTimeout != time.Second {
		t.Fatalf("expected read header timeout 1s, got %s", server.server.ReadHeaderTimeout)
	}
	if server.server.WriteTimeout != 2*time.Second {
		t.Fatalf("expected write timeout 2s, got %s", server.server.WriteTimeout)
	}
	if server.server.IdleTimeout != 3*time.Second {
		t.Fatalf("expected idle timeout 3s, got %s", server.server.IdleTimeout)
	}
}
