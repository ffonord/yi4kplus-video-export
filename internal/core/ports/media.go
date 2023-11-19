package ports

import (
	"context"
	"github.com/ffonord/yi4kplus-video-export/internal/core/domain/file"
	"io"
)

type Media interface {
	SessionStart(ctx context.Context) error
	GetFiles(ctx context.Context) (<-chan *file.File, error)
	GetReader(f *file.File) (io.ReadCloser, error)
	Delete(f *file.File) error
}
