package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/robatipoor/task-scheduler/internal/master/config"
)

type WorkerClient struct {
	client *http.Client
}

type WorkerClientInterface interface {
	HealthCheck(baseUrl string) (*int, error)
	AssignTask(baseUrl string, reqBody AssingTaskWorkerRequest) (*int, error)
	GetResultTask(baseUrl string, id uint) (*ResultTaskResponse, error)
}

func NewWorkerClient(cfg *config.Configure) *WorkerClient {
	return &WorkerClient{client: &http.Client{
		Timeout: cfg.Client.Timeout,
	}}
}

type AssingTaskWorkerRequest struct {
	TrackID uint `json:"track_id"`
	Input   uint `json:"input"`
}

type ResultTaskResponse struct {
	Status  *uint `json:"status"`
	Result *uint `json:"result"`
}

func (sc *WorkerClient) GetResultTask(baseUrl string, trackID uint) (*ResultTaskResponse, error) {
	url := fmt.Sprintf("%s/api/v1/tasks/result/%d", baseUrl, trackID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Error closing response body: ", err)
		}
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response ResultTaskResponse
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	return &response, nil
}

func (sc *WorkerClient) AssignTask(baseUrl string, reqBody AssingTaskWorkerRequest) (*int, error) {
	url := fmt.Sprintf("%s/api/v1/tasks/submit", baseUrl)
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Error closing response body: ", err)
		}
	}()

	return &resp.StatusCode, nil
}

func (sc *WorkerClient) HealthCheck(baseUrl string) (*int, error) {
	url := fmt.Sprintf("%s/health", baseUrl)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := sc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Println("Error closing response body: ", err)
		}
	}()

	return &resp.StatusCode, nil
}
