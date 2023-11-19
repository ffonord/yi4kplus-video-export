package app

import (
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/amba"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/ftp"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/telnet"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/storage/localdisk"
	"github.com/ffonord/yi4kplus-video-export/internal/core/services/mediaexporter"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
)

type App struct {
	MediaExporter *mediaexporter.MediaExporter
	Logger        *logger.Logger
}

func New(env string) *App {
	log := logger.New(env)

	ambaTCPConnFactory := new(amba.NetTCPConnFactory)
	ambaConfig := amba.NewConfig()
	ambaClient := amba.New(ambaConfig, log, ambaTCPConnFactory)

	telnetTCPConnFactory := new(telnet.NetTCPConnFactory)
	telnetBufioReaderFactory := new(telnet.BufioReaderFactory)
	telnetConfig := telnet.NewConfig()
	telnetClient := telnet.New(telnetConfig, log, telnetTCPConnFactory, telnetBufioReaderFactory)

	ftpConfig := ftp.NewConfig()
	ftpConnFactory := new(ftp.FTPConnFactory)
	ftpClient := ftp.New(ftpConfig, log, ftpConnFactory)

	mediaDevice := yi4kplus.New(ambaClient, ftpClient, telnetClient)

	storageConfig := localdisk.NewConfig()
	storage := localdisk.New(storageConfig, log)

	me := mediaexporter.New(mediaDevice, storage, log)

	return &App{
		MediaExporter: me,
		Logger:        log,
	}
}
