package dto

import "github.com/robatipoor/task-scheduler/internal/worker/models"

type RegisterResponse struct {
	ID uint `json:"id"`
}

type TaskResponse struct {
	ID uint `json:"id"`
}

type TaskResultResponse struct {
	Result *uint `json:"result"`
	Status *models.TaskStatus  `json:"status"`
}

type HealthCheckResponse struct {
	Services []ServiceCheckStatus `json:"services"`
}

type ServiceCheckStatus struct {
	ServiceName string `json:"service_name"`
	IsReady     bool   `json:"is_ready"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
