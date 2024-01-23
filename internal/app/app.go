package app

import (
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/amba"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/ftp"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/telnet"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/storage/localdisk"
	"github.com/ffonord/yi4kplus-video-export/internal/core/ports"
	"github.com/ffonord/yi4kplus-video-export/internal/core/services/filehandler"
	"github.com/ffonord/yi4kplus-video-export/internal/core/services/mediaexporter"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
)

type App struct {
	MediaExporter *mediaexporter.MediaExporter
	FileHandler   *filehandler.FileHandler
	Logger        *logger.Logger
}

func New(env string) *App {
	log := logger.New(env)

	ambaTCPConnFactory := new(amba.NetTCPConnFactory)
	ambaBufioReaderFactory := new(amba.BufioReaderFactory)
	ambaConfig := amba.NewConfig()
	ambaClient := amba.New(ambaConfig, log, ambaTCPConnFactory, ambaBufioReaderFactory)

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

	encoders := map[string]ports.Encoder{}

	fileHandlerConfig := filehandler.NewConfig()
	fh := filehandler.New(fileHandlerConfig, log, encoders)

	return &App{
		MediaExporter: me,
		FileHandler:   fh,
		Logger:        log,
	}
}
