package config

import (
	"os"

	"github.com/rs/zerolog"
)

func LogLevel() string {
	if v := os.Getenv("LOG_LEVEL"); v != "" {
		return v
	}
	return zerolog.InfoLevel.String()
}

func HostPort() string {
	if v := os.Getenv("HOST_PORT"); v != "" {
		return v
	}
	return ":8000"
}
