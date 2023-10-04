package main

import "github.com/ffonord/yi4kplus-video-export/internal"

func main() {
	client := internal.NewClient("yi4k", "7878")

	client.Run()
	client.Stop()
}
