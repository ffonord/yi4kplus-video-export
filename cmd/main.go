package main

import (
	"context"
	"errors"
	"flag"
	"github.com/ffonord/yi4kplus-video-export/internal/app"
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
	cancelTimeout       = time.Second * 3
)

var (
	envFilePath string
)

func init() {
	flag.StringVar(&envFilePath, "env-file-path", ".env", "path to .env file with variables")
}

// TODO: добавить команду декодирования через ffmpeg
// пример запуска из golang приложения: https://gist.github.com/tpanum/374ca0a415a5d94ffcd9
// пример команды `ffmpeg -i Y0030857.MP4 -vcodec libx265 -crf 28 output.mp4` - размер файла уменьшился в ~20 раз
func main() {
	handleError(loadDotEnvFile(), "gotenv load file")

	apl := app.New(os.Getenv("ENV"))
	var wg sync.WaitGroup

	ctx, cancelFunc := context.WithCancel(context.Background())
	wg.Add(1)

	go func() {
		defer wg.Done()

		err := apl.MediaExporter.Run(ctx, surveyPeriod, autoShutdownTimeOut)
		apl.Logger.Error(errors.Unwrap(err).Error())
	}()

	wait(&wg, cancelFunc)
}

func loadDotEnvFile() error {
	flag.Parse()
	return gotenv.Load(envFilePath)
}

func wait(wg *sync.WaitGroup, cancelFunc context.CancelFunc) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	<-s

	cancelFunc()

	doneChan := make(chan struct{})
	go func() {
		wg.Wait()
		close(doneChan)
	}()

	select {
	case <-time.After(cancelTimeout):
	case <-doneChan:
	}
}

func handleError(e error, message string) {
	if e != nil {
		log.Fatalf("%s error: %s", message, e.Error())
	}
}
