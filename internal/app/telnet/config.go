package telnet

import "os"

type Config struct {
	host          string
	port          string
	user          string
	ftpServerPort string
	ftpServerUser string
	mediaDir      string
	logLevel      string
}

func NewConfig() *Config {
	return &Config{
		host:          os.Getenv("TELNET_SERVER_HOST"),
		port:          os.Getenv("TELNET_SERVER_PORT"),
		user:          os.Getenv("TELNET_SERVER_USER"),
		ftpServerPort: os.Getenv("FTP_SERVER_PORT"),
		ftpServerUser: os.Getenv("FTP_SERVER_USER"),
		mediaDir:      os.Getenv("SERVER_MEDIA_DIR"),
		logLevel:      os.Getenv("TELNET_SERVER_LOG_LEVEL"),
	}
}
