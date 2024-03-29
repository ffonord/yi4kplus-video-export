package ftp

import "os"

type Config struct {
	host     string
	port     string
	user     string
	password string
}

func NewConfig() *Config {
	return &Config{
		host: os.Getenv("FTP_SERVER_HOST"),
		port: os.Getenv("FTP_SERVER_PORT"),
		user: os.Getenv("FTP_SERVER_USER"),
	}
}
