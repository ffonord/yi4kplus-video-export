package main

import (
	"context"
	"github.com/ffonord/yi4kplus-video-export/internal"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client := internal.NewClient("yi4k", "7878")

	client.Run()
	client.Stop()
}
