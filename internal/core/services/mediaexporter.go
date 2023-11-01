package services

import (
	"context"
	"fmt"
	"github.com/ffonord/yi4kplus-video-export/internal/core/ports"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	"io"
)

type MediaExporter struct {
	mediaAdapter   ports.Media
	storageAdapter ports.Storage
	logger         *logger.Logger
}

func NewMediaExporter(media ports.Media, storage ports.Storage) *MediaExporter {
	return &MediaExporter{
		mediaAdapter:   media,
		storageAdapter: storage,
		logger:         logger.New(),
	}
}

func (e *MediaExporter) ExportFiles(ctx context.Context) error {
	ffCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := e.mediaAdapter.SessionStart(ffCtx)
	if err != nil {
		return e.errWrap("ExportFiles", "media adapter session start", err)
	}

	err = e.storageAdapter.SessionStart(ffCtx)
	if err != nil {
		return e.errWrap("ExportFiles", "storage adapter session start", err)
	}

	fileChan, err := e.mediaAdapter.GetFiles(ctx)
	if err != nil {
		return e.errWrap("ExportFiles", "media adapter get files", err)
	}

	for f := range fileChan {

		dstFileWriter, err := e.storageAdapter.GetWriter(f)
		if err != nil {
			return e.errWrap("ExportFiles", "storage adapter get writer", err)
		}

		srcFileReader, err := e.mediaAdapter.GetReader(f)
		if err != nil {
			return e.errWrap("ExportFiles", "media adapter get reader", err)
		}

		e.logger.Infof("Start download %s", f.Name)

		written, err := io.Copy(dstFileWriter, srcFileReader)
		if err != nil {
			return e.errWrap("ExportFiles", "io copy "+f.Path, err)
		}

		err = srcFileReader.Close()
		if err != nil {
			return e.errWrap("ExportFiles", "media adapter reader close", err)
		}

		err = dstFileWriter.Close()
		if err != nil {
			return e.errWrap("ExportFiles", "storage adapter writer close", err)
		}

		if uint64(written) == f.Size {
			err := e.mediaAdapter.Delete(f)
			if err != nil {
				return e.errWrap("ExportFiles", "media adapter delete", err)
			}

			e.logger.Infof("Success downloaded %s", f.Name)
		} else {
			err := e.storageAdapter.Delete(f)
			if err != nil {
				return e.errWrap("ExportFiles", "storage adapter delete", err)
			}
		}
	}

	return nil
}

func (e *MediaExporter) errWrap(methodName, message string, err error) error {
	return fmt.Errorf("\n\tmediaexporter::%s: %s failed: %w\n", methodName, message, err)
}
