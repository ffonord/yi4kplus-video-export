package local

import (
	"context"
	"fmt"
	"github.com/ffonord/yi4kplus-video-export/internal/core/domain"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	"io"
)

type Storage struct {
	config *Config
	logger *logger.Logger
}

func NewStorage(config *Config) *Storage {
	return &Storage{
		config: config,
		logger: logger.New(),
	}
}

func (s *Storage) configureLogger() error {
	return s.logger.SetLevel(s.config.logLevel)
}

func (s *Storage) SessionStart(ctx context.Context) error {
	err := s.configureLogger()
	if err != nil {
		return s.errWrap("Run", "configure logger", err)
	}

	//TODO:

	return nil
}

func (s *Storage) GetWriter(f *domain.File) (io.WriteCloser, error) {
	//TODO:
	return nil, nil
}

func (s *Storage) Delete(f *domain.File) error {
	//TODO:
	return nil
}

func (s *Storage) errWrap(methodName, message string, err error) error {
	return fmt.Errorf("\n\tlocalStorage::%s: %s failed: %w\n", methodName, message, err)
}
