package ports

import (
	"context"
	"github.com/ffonord/yi4kplus-video-export/internal/core/domain"
	"io"
)

type Storage interface {
	SessionStart(ctx context.Context) error
	GetWriter(f *domain.File) (io.WriteCloser, error)
	Delete(f *domain.File) error
}
