package ports

import (
	"context"
	"github.com/ffonord/yi4kplus-video-export/internal/core/domain"
	"io"
)

type Media interface {
	SessionStart(ctx context.Context) error
	GetFiles(ctx context.Context) (<-chan *domain.File, error)
	GetReader(f *domain.File) (io.ReadCloser, error)
	Delete(f *domain.File) error
}
