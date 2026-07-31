// Package config initializes and exposes the application's startup configuration.
package config

import (
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// Config contains the process settings read once during startup.
type Config struct {
	HTTP HTTPConfig
	Log  LogConfig
}

// HTTPConfig configures the HTTP listener.
type HTTPConfig struct {
	Address           string
	ReadHeaderTimeout time.Duration
	WriteTimeout      time.Duration
	IdleTimeout       time.Duration
	ShutdownTimeout   time.Duration
}

// LogConfig configures process-wide logging.
type LogConfig struct {
	Level string
}

var (
	cfgOnce sync.Once
	config  *Config
)

// Initialize sets configuration defaults, reads environment variables through Viper,
// and initializes the process-wide configuration singleton. The first call wins.
func Initialize() {
	cfgOnce.Do(func() {
		v := viper.New()
		setDefaults(v)
		v.AutomaticEnv()

		config = &Config{
			HTTP: HTTPConfig{
				Address:           strings.TrimSpace(v.GetString("HTTP_ADDRESS")),
				ReadHeaderTimeout: v.GetDuration("HTTP_READ_HEADER_TIMEOUT"),
				WriteTimeout:      v.GetDuration("HTTP_WRITE_TIMEOUT"),
				IdleTimeout:       v.GetDuration("HTTP_IDLE_TIMEOUT"),
				ShutdownTimeout:   v.GetDuration("HTTP_SHUTDOWN_TIMEOUT"),
			},
			Log: LogConfig{
				Level: strings.TrimSpace(v.GetString("LOG_LEVEL")),
			},
		}
	})
}

func setDefaults(v *viper.Viper) {
	v.SetDefault("HTTP_ADDRESS", ":8083")
	v.SetDefault("HTTP_READ_HEADER_TIMEOUT", 5*time.Second)
	v.SetDefault("HTTP_WRITE_TIMEOUT", 10*time.Second)
	v.SetDefault("HTTP_IDLE_TIMEOUT", 60*time.Second)
	v.SetDefault("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second)
	v.SetDefault("LOG_LEVEL", "info")
}

// GetConfig returns the process-wide configuration, initializing it if needed.
func GetConfig() *Config {
	Initialize()
	return config
}
