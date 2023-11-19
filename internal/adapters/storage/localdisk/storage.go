package localdisk

import (
	"context"
	"errors"
	"fmt"
	"github.com/ffonord/yi4kplus-video-export/internal/core/domain/file"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	"io"
	"log/slog"
	"os"
	"syscall"
)

var (
	ErrNotEnoughSpace = errors.New("the disk has run out of free space")
)

type Storage struct {
	config *Config
	logger *logger.Logger
}

func New(config *Config, logger *logger.Logger) *Storage {
	return &Storage{
		config: config,
		logger: logger,
	}
}

func (s *Storage) checkFreeMemory() error {
	const op = "Storage.checkFreeMemory"

	log := s.logger.With(
		slog.String("op", op),
		slog.Any("config", s.config),
	)

	fs := syscall.Statfs_t{}
	err := syscall.Statfs(s.config.storageDir, &fs)
	if err != nil {
		return s.errWrap(op, "syscall statfs", err)
	}

	availableBytes := fs.Bfree * uint64(fs.Bsize)
	availableGigaBytes := float64(availableBytes) / float64(1024*1024*1024)

	log.Info(fmt.Sprintf("Free: %.2f Gb\n", availableGigaBytes))

	//TODO: добавить в конфиг "10"
	if availableBytes < 10 {
		return ErrNotEnoughSpace
	}

	return nil
}

func (s *Storage) SessionStart(ctx context.Context) error {
	return s.checkFreeMemory()
}

func (s *Storage) GetWriter(f *file.File) (io.WriteCloser, error) {
	const op = "Storage.GetWriter"

	file, err := os.Open(f.Path)

	if os.IsNotExist(err) {
		file, err = os.Create(f.Path)
		if err != nil {
			return nil, s.errWrap(op, "os create", err)
		}
	} else if err != nil {
		return nil, s.errWrap(op, "os open", err)
	}

	return file, nil
}

func (s *Storage) Delete(f *file.File) error {
	const op = "Storage.Delete"

	err := os.Remove(f.Path)
	if err != nil {
		return s.errWrap(op, "os remove", err)
	}

	return nil
}

func (s *Storage) errWrap(methodName, message string, err error) error {
	return fmt.Errorf("%s: %s failed: %w", methodName, message, err)
}
