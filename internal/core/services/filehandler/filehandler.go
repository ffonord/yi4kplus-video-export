package filehandler

import (
	"context"
	"fmt"
	"github.com/ffonord/yi4kplus-video-export/internal/core/ports"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	"log/slog"
	"sync"
	"time"
)

type FileHandler struct {
	config   *Config
	logger   *logger.Logger
	encoders map[string]ports.Encoder
}

func New(config *Config, logger *logger.Logger, encoders map[string]ports.Encoder) *FileHandler {
	return &FileHandler{
		config:   config,
		logger:   logger,
		encoders: encoders,
	}
}

func (fe *FileHandler) Run(ctx context.Context) {
	const op = "FileHandler.Run"

	log := fe.logger.With(
		slog.String("op", op),
	)

	wg := &sync.WaitGroup{}

	for {
		for name, encoder := range fe.encoders {
			wg.Add(1)

			go func(encoder ports.Encoder, name string) {
				defer wg.Done()

				err := encoder.Encode(ctx, fe.config.storageDir, fe.config.encodedDir)

				if err != nil {
					log.Error(fmt.Errorf("failed encode file with %s-encoder, err: %w", name, err).Error())
				}
			}(encoder, name)
		}

		wg.Wait()

		select {
		case <-ctx.Done():
			log.Info(fmt.Sprintf("Success stop %s", op))
			return
		case <-time.After(fe.config.pollingMinutes):
		}
	}
}
