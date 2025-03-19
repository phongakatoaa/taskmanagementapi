package postgres

import (
	"fmt"
	"os"
)

const (
	hostKey     = "PG_HOST"
	portKey     = "PG_PORT"
	usernameKey = "PG_USERNAME"
	passwordKey = "PG_PASSWORD"
	databaseKey = "PG_DATABASE"
	sslModeKey  = "PG_SSLMODE"
)

func host() string {
	if v := os.Getenv(hostKey); v != "" {
		return v
	}
	return "localhost"
}

func port() string {
	if v := os.Getenv(portKey); v != "" {
		return v
	}
	return "5432"
}

func username() string {
	if v := os.Getenv(usernameKey); v != "" {
		return v
	}
	return "postgres"
}

func password() string {
	if v := os.Getenv(passwordKey); v != "" {
		return v
	}
	return "password"
}

func sslmode() string {
	if v := os.Getenv(sslModeKey); v != "" {
		return v
	}
	return "disable"
}

func database() string {
	if v := os.Getenv(databaseKey); v != "" {
		return v
	}
	return "task_management"
}

func ConnectionString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=%s dbname=%s",
		host(), port(), username(), password(), sslmode(), database())
}
