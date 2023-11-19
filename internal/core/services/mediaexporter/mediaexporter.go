package mediaexporter

import (
	"context"
	"errors"
	"fmt"
	"github.com/ffonord/yi4kplus-video-export/internal/core/ports"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	"io"
	"log/slog"
	"time"
)

type MediaExporter struct {
	mediaAdapter   ports.Media
	storageAdapter ports.Storage
	logger         *logger.Logger
}

func New(media ports.Media, storage ports.Storage, logger *logger.Logger) *MediaExporter {
	return &MediaExporter{
		mediaAdapter:   media,
		storageAdapter: storage,
		logger:         logger,
	}
}

func (e *MediaExporter) Run(ctx context.Context, surveyPeriod, delayPeriod time.Duration) error {
	const op = "MediaExporter.Run"

	log := e.logger.With(
		slog.String("op", op),
	)

	for {

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		err := e.ExportFiles(ctx)
		waitTime := surveyPeriod

		if err == nil {
			waitTime = delayPeriod
		} else {
			log.Info(errors.Unwrap(err).Error())
		}

		time.Sleep(waitTime)
	}
}

func (e *MediaExporter) ExportFiles(ctx context.Context) error {
	const op = "MediaExporter.ExportFiles"

	log := e.logger.With(
		slog.String("op", op),
	)

	ffCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := e.mediaAdapter.SessionStart(ffCtx)
	if err != nil {
		return e.errWrap(op, "media adapter session start", err)
	}

	err = e.storageAdapter.SessionStart(ffCtx)
	if err != nil {
		return e.errWrap(op, "storage adapter session start", err)
	}

	fileChan, err := e.mediaAdapter.GetFiles(ffCtx)
	if err != nil {
		return e.errWrap(op, "media adapter get files", err)
	}

	for f := range fileChan {

		dstFileWriter, err := e.storageAdapter.GetWriter(f)
		if err != nil {
			return e.errWrap(op, "storage adapter get writer", err)
		}

		srcFileReader, err := e.mediaAdapter.GetReader(f)
		if err != nil {
			return e.errWrap(op, "media adapter get reader", err)
		}

		log.Info("Start download: " + f.Name)

		written, err := io.Copy(dstFileWriter, srcFileReader)
		if err != nil {
			return e.errWrap(op, "io copy "+f.Path, err)
		}

		err = srcFileReader.Close()
		if err != nil {
			return e.errWrap(op, "media adapter reader close", err)
		}

		err = dstFileWriter.Close()
		if err != nil {
			return e.errWrap(op, "storage adapter writer close", err)
		}

		if uint64(written) == f.Size {
			err := e.mediaAdapter.Delete(f)
			if err != nil {
				return e.errWrap(op, "media adapter delete", err)
			}

			log.Info("Success downloaded: " + f.Name)
		} else {
			err := e.storageAdapter.Delete(f)
			if err != nil {
				return e.errWrap(op, "storage adapter delete", err)
			}
		}
	}

	return nil
}

func (e *MediaExporter) errWrap(methodName, message string, err error) error {
	return fmt.Errorf("%s: %s failed: %w", methodName, message, err)
}
