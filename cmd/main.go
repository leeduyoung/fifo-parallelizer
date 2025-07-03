package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/leeduyoung/fifo-parallelizer/internal/container"
)

func main() {
	app, err := container.NewContainer()
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := app.WorkerPool().Start(ctx); err != nil {
			log.Printf("Worker pool error: %v", err)
		}
	}()

	<-sigChan
	log.Println("Received signal to stop")

	cancel()
}
