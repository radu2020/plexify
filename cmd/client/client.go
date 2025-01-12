package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const serverURL = "http://localhost:8080"

type JobResponse struct {
	JobID int `json:"job_id"`
}

type StatusResponse struct {
	JobID  int    `json:"job_id"`
	Status string `json:"status"`
}

func main() {
	clientCount := 2
	requestsPerClient := 2
	periodicChecks := 10
	wg := sync.WaitGroup{}

	// Simulate multiple clients
	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		go simulateClient(i, requestsPerClient, periodicChecks, &wg)
	}

	wg.Wait()
	log.Println("All clients finished.")
}

// Simulate a single client
func simulateClient(clientID, requests, periodicChecks int, wg *sync.WaitGroup) {
	defer wg.Done()
	jobIDs := make([]int, 0)

	// Send multiple POST /job requests
	log.Printf("Client %d: Sending %d job requests...\n", clientID, requests)
	for i := 0; i < requests; i++ {
		payload := map[string]string{"payload": fmt.Sprintf("Client %d Job %d", clientID, i)}
		jobID := sendJobRequest(payload)
		jobIDs = append(jobIDs, jobID)
		time.Sleep(time.Millisecond * 500) // Simulate delay between job submissions
	}
	
	// Send multiple GET /status/{job_id} requests
	// Each periodic check, checks the status for all existing jobs.
	start := time.Now()
	for i := 0; i < periodicChecks; i++ {
		log.Printf("Client %d: Querying job status start...\n", clientID)
		
		for _, jobID := range jobIDs {
			queryJobStatus(clientID, jobID)
			time.Sleep(time.Millisecond * 5000) // Simulate delay between status checks
		}

		t := time.Now()
		elapsed := t.Sub(start)

		log.Printf("Client %d: Querying job status end... %s elapsed", clientID, elapsed)
	}

}

// Send a POST /job request
func sendJobRequest(payload map[string]string) int {
	body, _ := json.Marshal(payload)
	resp, err := http.Post(serverURL+"/job", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error sending job request: %v", err)
		return -1
	}
	defer resp.Body.Close()

	var jobResp JobResponse
	if err := json.NewDecoder(resp.Body).Decode(&jobResp); err != nil {
		log.Printf("Error decoding job response: %v", err)
		return -1
	}

	log.Printf("Job created: ID = %d", jobResp.JobID)
	return jobResp.JobID
}

// Query a job's status using GET /status/{job_id}
func queryJobStatus(clientID, jobID int) {
	resp, err := http.Get(fmt.Sprintf("%s/status/%d", serverURL, jobID))
	if err != nil {
		log.Printf("Client %d: Error querying job %d status: %v", clientID, jobID, err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Client %d: Status for job %d: %s", clientID, jobID, body)
}
