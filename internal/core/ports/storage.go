package ports

import (
	"context"
	"github.com/ffonord/yi4kplus-video-export/internal/core/domain/file"
	"io"
)

type Storage interface {
	SessionStart(ctx context.Context) error
	GetWriter(f *file.File) (io.WriteCloser, error)
	Delete(f *file.File) error
}
