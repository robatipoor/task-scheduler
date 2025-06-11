package services

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/robatipoor/task-scheduler/internal/worker/client"
	"github.com/robatipoor/task-scheduler/internal/worker/config"
	"github.com/robatipoor/task-scheduler/internal/worker/dto"
	"github.com/robatipoor/task-scheduler/internal/worker/models"
	"github.com/robatipoor/task-scheduler/internal/worker/repositories"
	"gorm.io/gorm"
)

type TaskService struct {
	taskRepo     repositories.TaskRepositoryInterface
	context      context.Context
	wg           sync.WaitGroup
	config       *config.Configure
	masterClient client.MasterClientInterface
}

type TaskServiceInterface interface {
	CheckExistTask(taskReq dto.SubmitTaskRequest) (*models.Task, error)
	SubmitTask(taskReq dto.SubmitTaskRequest) (*models.Task, error)
	GetResultTask(TrackUID uint) (*models.Task, error)
	Wait()
}

func NewTaskService(
	context context.Context,
	config *config.Configure,
	masterClient client.MasterClientInterface,
	taskRepo repositories.TaskRepositoryInterface,
) *TaskService {
	return &TaskService{
		context:      context,
		config:       config,
		taskRepo:     taskRepo,
		masterClient: masterClient,
	}
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
	task, err := ts.taskRepo.Save(taskReq.TrackID, taskReq.Input)

	if err != nil {
		return nil, err
	}

	go func() {
		defer ts.wg.Done()

		select {
		case <-time.After(time.Duration(task.Input) * time.Millisecond):
			ts.wg.Add(1)
			r := uint(rand.Intn(int(task.Input)))

			if _, err := ts.taskRepo.Update(task.TrackID, r, models.Reporting); err != nil {
				log.Printf("Failed to update task %v", err)
				return
			}

			req := client.UpdateResultTaskRequest{
				TrackID: task.TrackID,
				Result:  r,
			}

			statusCode, err := ts.masterClient.UpdateResultTask(ts.config.Master.Url, req)
			if err != nil {
				log.Printf("Failed to update result task: %v", err)
				return
			}

			if *statusCode != 200 {
				log.Println("Failed to update result task: ", statusCode)
				return
			}

			if _, err := ts.taskRepo.Update(task.TrackID, r, models.Completed); err != nil {
				log.Printf("Failed to update task %v", err)
				return
			}

		case <-ts.context.Done():
			log.Println("Goroutine cancelled:", context.Canceled.Error())
			return
		}
	}()

	return task, nil
}

func (ts *TaskService) GetResultTask(TrackID uint) (*models.Task, error) {
	task, err := ts.taskRepo.FindByTrackID(TrackID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (ts *TaskService) Wait() {
	ts.wg.Wait()
}