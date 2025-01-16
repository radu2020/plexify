package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

var (
	JobIDCounter int32
	JobQueue     chan Job
	JobStatuses  sync.Map
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
type StringJobProcessor struct {
	JobStatuses  sync.Map
	JobQueue     chan Job
	JobIDCounter int32
}

// Process simulates processing jobs.
func (sjp *StringJobProcessor) Process(job Job) error {
	processTime := time.Duration(rand.Intn(26)+5) * time.Second
	time.Sleep(processTime)
	return nil
}

// HTTP handler to create a job
func (sjp *StringJobProcessor) CreateJobHandler(w http.ResponseWriter, r *http.Request) {
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

	jobID := int(atomic.AddInt32(&sjp.JobIDCounter, 1))
	job := Job{ID: jobID, Payload: payload.Payload}

	sjp.JobStatuses.Store(jobID, "pending")
	sjp.JobQueue <- job

	w.WriteHeader(http.StatusAccepted)
	response := map[string]int{"job_id": jobID}
	json.NewEncoder(w).Encode(response)
}

// HTTP handler to check job status
func (sjp *StringJobProcessor) JobStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var jobID int
	if _, err := fmt.Sscanf(r.URL.Path, "/status/%d", &jobID); err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	if status, ok := sjp.JobStatuses.Load(jobID); ok {
		response := map[string]interface{}{
			"job_id": jobID,
			"status": status,
		}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Job not found", http.StatusNotFound)
	}
}

// Worker function which processes jobs
func Worker(ctx context.Context, wg *sync.WaitGroup, id int, sjp *StringJobProcessor) {
	defer wg.Done()

	for {
		select {
		case job := <-sjp.JobQueue:
			log.Printf("Worker %d: Started job %d\n", id, job.ID)
			sjp.JobStatuses.Store(job.ID, "processing")
			if err := sjp.Process(job); err != nil {
				log.Printf("Worker %d: Failed to process job %d: %v\n", id, job.ID, err)
				sjp.JobStatuses.Store(job.ID, "failed")
				continue
			}
			log.Printf("Worker %d: Completed job %d\n", id, job.ID)
			sjp.JobStatuses.Store(job.ID, "completed")

		case <-ctx.Done():
			log.Printf("Worker %d: Received cancellation signal\n", id)
			return
		}
	}
}
