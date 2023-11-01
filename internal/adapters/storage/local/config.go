package local

import "os"

type Config struct {
	storageDir string
	logLevel   string
}

func NewConfig() *Config {
	return &Config{
		storageDir: os.Getenv("LOCAL_STORAGE_DIR"),
		logLevel:   os.Getenv("LOCAL_STORAGE_LOG_LEVEL"),
	}
}
