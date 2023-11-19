package localdisk

import "os"

type Config struct {
	storageDir string
}

func NewConfig() *Config {
	return &Config{
		storageDir: os.Getenv("LOCAL_STORAGE_DIR"),
	}
}
