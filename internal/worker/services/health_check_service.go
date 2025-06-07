package services

import (
	"log"

	"github.com/robatipoor/task-scheduler/internal/worker/dto"
	"gorm.io/gorm"
)

type HealthCheckService struct {
	db *gorm.DB
}

type HealthCheckServiceInterface interface {
	Check() *dto.HealthCheckResponse
}

func NewHealthCheckService(db *gorm.DB) *HealthCheckService {
	return &HealthCheckService{db: db}
}

func (s *HealthCheckService) Check() *dto.HealthCheckResponse {
	isReady := true
	if err := s.db.Exec("SELECT 1").Error; err != nil {
		log.Printf("postgres database is not healthy: %v \n", err)
		isReady = false
	}
	return &dto.HealthCheckResponse{
		Services: []dto.ServiceCheckStatus{{
			ServiceName: "postgres",
			IsReady:     isReady,
		}},
	}
}
