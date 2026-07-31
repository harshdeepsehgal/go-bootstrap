package config

import (
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestHTTPTimeoutDefaults(t *testing.T) {
	v := viper.New()
	setDefaults(v)

	if timeout := v.GetDuration("HTTP_WRITE_TIMEOUT"); timeout != 10*time.Second {
		t.Fatalf("expected write timeout 10s, got %s", timeout)
	}
	if timeout := v.GetDuration("HTTP_IDLE_TIMEOUT"); timeout != 60*time.Second {
		t.Fatalf("expected idle timeout 60s, got %s", timeout)
	}
}
