package main

import (
	"context"
	"github.com/radu2020/plexify/services/server/pkg/worker"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const workerCount = 5

func startServer(srv *http.Server) {
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe: %v", err)
	}
}

func main() {
	// Create Job Processor
	sjp := worker.StringJobProcessor{
		JobStatuses:  sync.Map{},
		JobQueue:     make(chan worker.Job, 100),
		JobIDCounter: 0,
	}

	// Channel to listen for system signals (e.g., Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Create a context that can be canceled
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx() // Ensure cancel is called at the end to clean up

	// Wait group to ensure all goroutines finish before exiting
	wg := sync.WaitGroup{}
	wg.Add(workerCount)

	// Start worker pool
	for id := 0; id < workerCount; id++ {
		go worker.Worker(ctx, &wg, id, &sjp)
	}

	// Register HTTP handlers
	http.HandleFunc("/job", sjp.CreateJobHandler)
	http.HandleFunc("/status/", sjp.JobStatusHandler)

	// Start HTTP Server
	srv := &http.Server{Addr: ":8080"}
	go startServer(srv)
	log.Println("Server started on :8080")

	// Wait for an interrupt signal to initiate graceful shutdown
	<-sigChan

	// Handle shutdown signal (Ctrl+C or SIGTERM)
	log.Println("Received shutdown signal. Shutting down gracefully...")

	// Cancel the context to notify all goroutines to stop
	cancelCtx()

	// Wait for all goroutines to finish
	wg.Wait()

	// Final cleanup before exiting
	log.Println("Server stopped gracefully.")
}
