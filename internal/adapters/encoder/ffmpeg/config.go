package ffmpeg

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	storageDir     string
	encodedDir     string
	pollingMinutes time.Duration
}

func NewConfig() *Config {
	pollingMinutes := time.Minute * 2
	rawPollingMinutes := os.Getenv("FILE_HANDLER_POLLING_MINUTES")

	if i, err := strconv.Atoi(rawPollingMinutes); err == nil {
		pollingMinutes = time.Minute * time.Duration(i)
	}

	return &Config{
		storageDir:     os.Getenv("LOCAL_STORAGE_DIR"),
		encodedDir:     os.Getenv("LOCAL_ENCODED_DIR"),
		pollingMinutes: pollingMinutes,
	}
}
