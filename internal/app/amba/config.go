package amba

import (
	"os"
)

type Config struct {
	host     string
	port     string
	logLevel string
}

func NewConfig() *Config {
	return &Config{
		host:     os.Getenv("AMBA_SERVER_HOST"),
		port:     os.Getenv("AMBA_SERVER_PORT"),
		logLevel: os.Getenv("AMBA_SERVER_LOG_LEVEL"),
	}
}
