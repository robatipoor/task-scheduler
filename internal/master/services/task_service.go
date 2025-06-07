package services

import (
	"github.com/robatipoor/task-scheduler/internal/master/dto"
	"github.com/robatipoor/task-scheduler/internal/master/models"
	"github.com/robatipoor/task-scheduler/internal/master/repositories"
	"gorm.io/gorm"
)

type TaskService struct {
	taskRepo repositories.TaskRepositoryInterface
}

type TaskServiceInterface interface {
	CheckExistTask(taskReq dto.SubmitTaskRequest) (*models.Task, error)
	SubmitTask(taskReq dto.SubmitTaskRequest) (*models.Task, error)
	GetResultTask(trackID string) (*models.AssignTask, error)
}

func NewTaskService(taskRepo repositories.TaskRepositoryInterface) *TaskService {
	return &TaskService{taskRepo: taskRepo}
}

func (ts *TaskService) CheckExistTask(taskReq dto.SubmitTaskRequest) (*models.Task, error) {
	task, err := ts.taskRepo.FindByTrackUID(taskReq.TrackUID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return task, nil
}

func (ts *TaskService) SubmitTask(taskReq dto.SubmitTaskRequest) (*models.Task, error) {
	return ts.taskRepo.Save(taskReq.TrackUID, taskReq.Input, taskReq.Priority)
}

func (ts *TaskService) GetResultTask(trackUID string) (*models.AssignTask, error) {
	task, err := ts.taskRepo.FindByTrackUID(trackUID)
	if err != nil {
		return nil, err
	}
	lt := len(task.AssignTasks)
	if lt >= 1 {
		for _, atask := range task.AssignTasks {
			if atask.Status == models.Completed {
				return &atask, nil
			}
		}
		return &task.AssignTasks[0], nil
	}
	return nil, nil
}
