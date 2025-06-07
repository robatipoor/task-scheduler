package services

import (
	"github.com/robatipoor/task-scheduler/internal/worker/dto"
	"github.com/robatipoor/task-scheduler/internal/worker/models"
	"github.com/robatipoor/task-scheduler/internal/worker/repositories"
	"gorm.io/gorm"
)

type TaskService struct {
	taskRepo repositories.TaskRepositoryInterface
}

type TaskServiceInterface interface {
	CheckExistTask(taskReq dto.SubmitTaskRequest) (*models.Task, error)
	SubmitTask(taskReq dto.SubmitTaskRequest) (*models.Task, error)
	GetResultTask(TrackUID uint) (*models.Task, error)
}

func NewTaskService(taskRepo repositories.TaskRepositoryInterface) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (ts *TaskService) CheckExistTask(taskReq dto.SubmitTaskRequest) (*models.Task, error) {
	task, err := ts.taskRepo.FindByTrackID(taskReq.TrackID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return task, nil
}

func (ts *TaskService) SubmitTask(taskReq dto.SubmitTaskRequest) (*models.Task, error) {
	return ts.taskRepo.Save(taskReq.TrackID, taskReq.Input)
}

func (ts *TaskService) GetResultTask(TrackID uint) (*models.Task, error) {
	task, err := ts.taskRepo.FindByTrackID(TrackID)
	if err != nil {
		return nil, err
	}
	return task, nil
}
