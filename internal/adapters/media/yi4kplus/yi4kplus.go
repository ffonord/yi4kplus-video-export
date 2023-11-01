package yi4kplus

import (
	"context"
	"fmt"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/amba"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/ftp"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/telnet"
	"github.com/ffonord/yi4kplus-video-export/internal/core/domain"
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

	err := y.ambaClient.Run(ctx)
	if err != nil {
		return y.errWrap("SessionStart", "amba client run", err)
	}

	err = y.telnetClient.Run(ctx)
	if err != nil {
		return y.errWrap("SessionStart", "telnet client run", err)
	}

	err = y.ftpClient.Run(ctx)
	if err != nil {
		return y.errWrap("SessionStart", "ftp client run", err)
	}

	return nil
}

func (y *Yi4kPlus) GetFiles(ctx context.Context) (<-chan *domain.File, error) {
	fileChan, err := y.ftpClient.GetFiles(ctx)
	if err != nil {
		return nil, y.errWrap("GetFiles", "ftp get files", err)
	}

	return fileChan, nil
}

func (y *Yi4kPlus) GetReader(f *domain.File) (io.ReadCloser, error) {
	reader, err := y.ftpClient.GetReader(f.Path)
	if err != nil {
		return nil, y.errWrap("GetReader", "ftp get reader "+f.Path, err)
	}
	return reader, nil
}

func (y *Yi4kPlus) Delete(f *domain.File) error {
	err := y.ftpClient.Delete(f.Path)
	if err != nil {
		return y.errWrap("Delete", "ftp delete", err)
	}

	return nil
}

func (y *Yi4kPlus) errWrap(methodName, message string, err error) error {
	return fmt.Errorf("\n\tyi4kplus::%s: %s failed: %w\n", methodName, message, err)
}
