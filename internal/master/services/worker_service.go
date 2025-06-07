package services

import (
	"github.com/robatipoor/task-scheduler/internal/master/dto"
	"github.com/robatipoor/task-scheduler/internal/master/models"
	"github.com/robatipoor/task-scheduler/internal/master/repositories"
)

type WorkerService struct {
	workerRepo     repositories.WorkerRepositoryInterface
	assingTaskRepo repositories.AssingTaskRepositoryInterface
}

type WorkerServiceInterface interface {
	Register(req dto.WorkerRegisterRequest) (*models.Worker, error)
	UpdateResultTask(req dto.WorkerUpdateResultTaskRequest) (*models.AssignTask, error)
}

func NewWorkerService(workerRepo repositories.WorkerRepositoryInterface, assingTaskRepo repositories.AssingTaskRepositoryInterface) *WorkerService {
	return &WorkerService{workerRepo: workerRepo, assingTaskRepo: assingTaskRepo}
}

func (ws *WorkerService) Register(req dto.WorkerRegisterRequest) (*models.Worker, error) {
	return ws.workerRepo.SaveOrUpdate(req.BaseUrl, models.Active)
}

func (ws *WorkerService) UpdateResultTask(req dto.WorkerUpdateResultTaskRequest) (*models.AssignTask, error) {
	return ws.assingTaskRepo.Update(req.TrackID, models.Completed, &req.Result, "")
}
