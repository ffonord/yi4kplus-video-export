package main

import (
	"context"
	"flag"
	"github.com/ffonord/yi4kplus-video-export/internal/app/amba"
	"github.com/ffonord/yi4kplus-video-export/internal/app/telnet"
	"github.com/ffonord/yi4kplus-video-export/internal/pkg/logger"
	"github.com/subosito/gotenv"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	shutdownTimeOut = 5 * time.Second
)

var (
	envFilePath string
)

func init() {
	flag.StringVar(&envFilePath, "env-file-path", ".env", "path to .env file with variables")
}

// cleanupFunc is a cleanup function on shutting down
type cleanupFunc func(ctx context.Context) error

// gracefulShutdown waits for termination syscalls and doing clean up operations after received it
// @see https://gist.github.com/aladhims/baea548df03be8f1a5f78a61636225f6
func gracefulShutdown(ctx context.Context, timeout time.Duration, log *logger.Logger, ops map[string]cleanupFunc) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscalls that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		log.Info("shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Infof("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for key, op := range ops {
			wg.Add(1)
			innerOp := op
			innerKey := key
			go func() {
				defer wg.Done()

				log.Infof("cleaning up: %s", innerKey)
				if err := innerOp(ctx); err != nil {
					log.Infof("%s: clean up failed: %s", innerKey, err.Error())
					return
				}

				log.Infof("%s was shutdown gracefully", innerKey)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}

func main() {
	flag.Parse()
	err := gotenv.Load(envFilePath)
	handleError(err, "gotenv load file")

	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	config := amba.NewConfig()
	client := amba.New(config)

	go func() {
		err := client.Run(ctx)
		handleError(err, "remote api client run")
	}()

	telnetConfig := telnet.NewConfig()
	telnetClient := telnet.New(telnetConfig)

	go func() {
		err := telnetClient.Run(ctx)
		handleError(err, "telnet client run")
	}()

	wait := gracefulShutdown(ctx, shutdownTimeOut, logger.New(), map[string]cleanupFunc{
		"amba-client":   client.Shutdown,
		"telnet-client": telnetClient.Shutdown,
	})

	<-wait
	os.Exit(0)
}

func handleError(e error, message string) {
	if e != nil {
		log.Fatalf("%s error: %s", message, e.Error())
	}
}
