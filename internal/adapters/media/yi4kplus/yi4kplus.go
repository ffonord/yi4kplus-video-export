package yi4kplus

import (
	"context"
	"fmt"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/amba"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/ftp"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/telnet"
	"github.com/ffonord/yi4kplus-video-export/internal/core/domain/file"
	"io"
)

type Yi4kPlus struct {
	ambaClient   *amba.Client
	ftpClient    *ftp.Client
	telnetClient *telnet.Client
}

func New(
	ambaClient *amba.Client,
	ftpClient *ftp.Client,
	telnetClient *telnet.Client,
) *Yi4kPlus {
	return &Yi4kPlus{
		ambaClient:   ambaClient,
		ftpClient:    ftpClient,
		telnetClient: telnetClient,
	}
}

func (y *Yi4kPlus) SessionStart(ctx context.Context) error {
	const op = "Yi4kPlus.SessionStart"

	err := y.ambaClient.Run(ctx)
	if err != nil {
		return y.errWrap(op, "amba client run", err)
	}

	err = y.telnetClient.Run(ctx)
	if err != nil {
		return y.errWrap(op, "telnet client run", err)
	}

	err = y.ftpClient.Run(ctx)
	if err != nil {
		return y.errWrap(op, "ftp client run", err)
	}

	return nil
}

func (y *Yi4kPlus) GetFiles(ctx context.Context) (<-chan *file.File, error) {
	const op = "Yi4kPlus.GetFiles"

	fileChan, err := y.ftpClient.GetFiles(ctx)
	if err != nil {
		return nil, y.errWrap(op, "ftp get files", err)
	}

	return fileChan, nil
}

func (y *Yi4kPlus) GetReader(f *file.File) (io.ReadCloser, error) {
	const op = "Yi4kPlus.GetReader"

	reader, err := y.ftpClient.GetReader(f.Path)
	if err != nil {
		return nil, y.errWrap(op, "ftp get reader "+f.Path, err)
	}
	return reader, nil
}

func (y *Yi4kPlus) Delete(f *file.File) error {
	const op = "Yi4kPlus.Delete"

	err := y.ftpClient.Delete(f.Path)
	if err != nil {
		return y.errWrap(op, "ftp delete "+f.Path, err)
	}

	return nil
}

func (y *Yi4kPlus) errWrap(methodName, message string, err error) error {
	return fmt.Errorf("%s: %s failed: %w", methodName, message, err)
}
