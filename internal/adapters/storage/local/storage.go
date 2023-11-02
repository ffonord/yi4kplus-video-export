package local

import (
	"context"
	"errors"
	"fmt"
	"github.com/ffonord/yi4kplus-video-export/internal/core/domain"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	"io"
	"os"
	"syscall"
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

	fs := syscall.Statfs_t{}
	err = syscall.Statfs(s.config.storageDir, &fs)
	if err != nil {
		return s.errWrap("SessionStart", "syscall statfs", err)
	}

	availableBytes := fs.Bfree * uint64(fs.Bsize)
	availableGigaBytes := float64(availableBytes) / float64(1024*1024*1024)

	s.logger.Infof("Free: %.2f Gb\n", availableGigaBytes)

	if availableBytes < 10 {
		return errors.New("the disk has run out of free space")
	}

	return nil
}

func (s *Storage) GetWriter(f *domain.File) (io.WriteCloser, error) {
	file, err := os.Open(f.Path)

	if os.IsNotExist(err) {
		file, err = os.Create(f.Path)
		if err != nil {
			return nil, s.errWrap("GetWriter", "os create", err)
		}
	} else if err != nil {
		return nil, s.errWrap("GetWriter", "os open", err)
	}

	return file, nil
}

func (s *Storage) Delete(f *domain.File) error {
	err := os.Remove(f.Path)
	if err != nil {
		return s.errWrap("Delete", "os remove", err)
	}

	return nil
}

func (s *Storage) errWrap(methodName, message string, err error) error {
	return fmt.Errorf("\n\tlocalStorage::%s: %s failed: %w\n", methodName, message, err)
}
