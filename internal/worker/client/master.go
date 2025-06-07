package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/robatipoor/task-scheduler/internal/worker/config"
)

type MasterClient struct {
	client *http.Client
}

type MasterClientInterface interface {
	Register(baseUrl string, reqBody WorkerRegisterRequest) (*int, error)
	UpdateResultTask(baseUrl string, reqBody UpdateResultTaskRequest) (*int, error)
}

func NewMasterClient(cfg *config.Configure) *MasterClient {
	return &MasterClient{client: &http.Client{
		Timeout: cfg.Client.Timeout,
	}}
}

type UpdateResultTaskRequest struct {
	TrackID uint `json:"track_id" binding:"required"`
	Result  uint `json:"result" binding:"required"`
}

type WorkerRegisterRequest struct {
	BaseUrl string `json:"base_url" binding:"required"`
}

func (mc *MasterClient) Register(baseUrl string, reqBody WorkerRegisterRequest) (*int, error) {
	url := fmt.Sprintf("%s/api/v1/workers/register", baseUrl)
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := mc.client.Do(req)
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

func (mc *MasterClient) UpdateResultTask(baseUrl string, reqBody UpdateResultTaskRequest) (*int, error) {
	url := fmt.Sprintf("%s/api/v1/workers/result", baseUrl)
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := mc.client.Do(req)
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
