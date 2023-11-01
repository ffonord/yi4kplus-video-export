package main

import (
	"context"
	"flag"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/amba"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/ftp"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/media/yi4kplus/telnet"
	"github.com/ffonord/yi4kplus-video-export/internal/adapters/storage/local"
	"github.com/ffonord/yi4kplus-video-export/internal/core/services"
	"github.com/subosito/gotenv"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	autoShutdownTimeOut = time.Minute * 5
	surveyPeriod        = time.Minute * 1
)

var (
	envFilePath string
)

func init() {
	flag.StringVar(&envFilePath, "env-file-path", ".env", "path to .env file with variables")
}

func main() {
	flag.Parse()
	err := gotenv.Load(envFilePath)
	handleError(err, "gotenv load file")

	me := initService()
	var vg sync.WaitGroup

	ctx, cancelFunc := context.WithCancel(context.Background())
	vg.Add(1)

	go func() {
		defer vg.Done()

		for {

			select {
			case <-ctx.Done():
				return
			default:
			}

			err = me.ExportFiles(ctx)
			waitTime := surveyPeriod

			if err == nil {
				waitTime = autoShutdownTimeOut
			}

			time.Sleep(waitTime)
		}
	}()

	wait(&vg, cancelFunc)
}

func initService() *services.MediaExporter {

	ambaConfig := amba.NewConfig()
	ambaClient := amba.New(ambaConfig)

	ftpConfig := ftp.NewConfig()
	ftpClient := ftp.New(ftpConfig)

	telnetConfig := telnet.NewConfig()
	telnetClient := telnet.New(telnetConfig)

	mediaDevice := yi4kplus.New(ambaClient, ftpClient, telnetClient)

	storageConfig := local.NewConfig()
	localStorage := local.NewStorage(storageConfig)

	return services.NewMediaExporter(mediaDevice, localStorage)
}

func wait(vg *sync.WaitGroup, cancelFunc context.CancelFunc) {

	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-s

	cancelFunc()

	doneChan := make(chan struct{})
	go func() {
		vg.Wait()
		close(doneChan)
	}()

	select {
	case <-time.After(time.Second * 10):
	case <-doneChan:
	}
}

func handleError(e error, message string) {
	if e != nil {
		log.Fatalf("%s error: %s", message, e.Error())
	}
}
