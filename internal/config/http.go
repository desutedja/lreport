package config

import "time"

type HTTPConfig struct {
	Port                    int
	GracefulShutdownTimeout time.Duration
}

func loadHTTPConfig() *HTTPConfig {
	return &HTTPConfig{
		Port:                    getIntWithDefault("HTTP_SERVER_PORT", 1234),
		GracefulShutdownTimeout: getDurationOrPanic("HTTP_GRACEFUL_SHUTDOWN_TIMEOUT"),
	}
}
