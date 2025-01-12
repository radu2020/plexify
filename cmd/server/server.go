package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"time"
)

type Job struct {
	ID      int
	Payload string
}

type JobStatus struct {
	ID     int
	Status string
}

// Interface for job processing
type JobProcessor interface {
	Process(job Job) error
}

// Concrete implementation for processing string-based jobs
type StringJobProcessor struct{}

func (sjp StringJobProcessor) Process(job Job) error {
	processTime := time.Duration(rand.Intn(26)+5) * time.Second
	time.Sleep(processTime) // Simulate job processing
	return nil
}

var (
	jobIDCounter int32
	jobQueue     chan Job
	jobStatuses  sync.Map
	workerCount  = 5
)

func main() {
	rand.Seed(time.Now().UnixNano())
	jobQueue = make(chan Job, 100)

	// Start worker pool
	for i := 0; i < workerCount; i++ {
		go worker(i)
	}

	http.HandleFunc("/job", createJobHandler)
	http.HandleFunc("/status/", jobStatusHandler)

	srv := &http.Server{Addr: ":8080"}

	// Graceful shutdown handling
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		log.Println("Shutting down server...")
		close(jobQueue) // Close job queue to stop workers
		srv.Close()
	}()

	log.Println("Server started on :8080")
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe: %v", err)
	}

	log.Println("Server stopped gracefully.")
}

// HTTP handler to create a job
func createJobHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		Payload string `json:"payload"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	jobID := int(atomic.AddInt32(&jobIDCounter, 1))
	job := Job{ID: jobID, Payload: payload.Payload}

	jobStatuses.Store(jobID, "pending")
	jobQueue <- job

	w.WriteHeader(http.StatusAccepted)
	response := map[string]int{"job_id": jobID}
	json.NewEncoder(w).Encode(response)
}

// HTTP handler to check job status
func jobStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var jobID int
	if _, err := fmt.Sscanf(r.URL.Path, "/status/%d", &jobID); err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	if status, ok := jobStatuses.Load(jobID); ok {
		response := map[string]interface{}{
			"job_id": jobID,
			"status": status,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Job not found", http.StatusNotFound)
	}
}

// Worker function to process jobs
func worker(id int) {
	processor := StringJobProcessor{}

	for job := range jobQueue {
		log.Printf("Worker %d: Started job %d\n", id, job.ID)
		jobStatuses.Store(job.ID, "processing")

		if err := processor.Process(job); err != nil {
			log.Printf("Worker %d: Failed to process job %d: %v\n", id, job.ID, err)
			jobStatuses.Store(job.ID, "failed")
			continue
		}

		log.Printf("Worker %d: Completed job %d\n", id, job.ID)
		jobStatuses.Store(job.ID, "completed")
	}
}
