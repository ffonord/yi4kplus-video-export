package amba

import (
	"os"
	"strconv"
)

type Config struct {
	host                string
	port                string
	autoShutdownTimeout int
}

func NewConfig() *Config {
	autoShutdownTimeout := 0
	rawAutoShutdownTimeout := os.Getenv("AMBA_SERVER_AUTO_SHUTDOWN_WITHOUT_CONNECTION_TIMEOUT")

	if i, err := strconv.Atoi(rawAutoShutdownTimeout); err == nil {
		autoShutdownTimeout = i
	}

	return &Config{
		host:                os.Getenv("AMBA_SERVER_HOST"),
		port:                os.Getenv("AMBA_SERVER_PORT"),
		autoShutdownTimeout: autoShutdownTimeout,
	}
}
