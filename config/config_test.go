package config

import (
	"testing"
	"time"

	"github.com/spf13/viper"
)

func TestSetDefaults(t *testing.T) {
	v := viper.New()
	setDefaults(v)

	tests := []struct {
		name string
		key  string
		want time.Duration
	}{
		{
			name: "write timeout",
			key:  "HTTP_WRITE_TIMEOUT",
			want: 10 * time.Second,
		},
		{
			name: "idle timeout",
			key:  "HTTP_IDLE_TIMEOUT",
			want: 60 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := v.GetDuration(tt.key); got != tt.want {
				t.Fatalf("expected %s, got %s", tt.want, got)
			}
		})
	}
}
