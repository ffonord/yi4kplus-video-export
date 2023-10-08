package telnet

import "os"

type Config struct {
	host     string
	port     string
	user     string
	logLevel string
}

func NewConfig() *Config {
	return &Config{
		host:     os.Getenv("TELNET_SERVER_HOST"),
		port:     os.Getenv("TELNET_SERVER_PORT"),
		user:     os.Getenv("TELNET_SERVER_USER"),
		logLevel: os.Getenv("TELNET_SERVER_LOG_LEVEL"),
	}
}
