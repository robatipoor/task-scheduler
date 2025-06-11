package services

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/robatipoor/task-scheduler/internal/worker/client"
	"github.com/robatipoor/task-scheduler/internal/worker/config"
	"github.com/robatipoor/task-scheduler/internal/worker/models"
	"github.com/robatipoor/task-scheduler/internal/worker/repositories"
)

type SchedulerService struct {
	context      context.Context
	config       *config.Configure
	wg           sync.WaitGroup
	taskRepo     repositories.TaskRepositoryInterface
	masterClient client.MasterClientInterface
}

func NewSchedulerService(
	context context.Context,
	config *config.Configure,
	masterClient client.MasterClientInterface,
	taskRepo repositories.TaskRepositoryInterface,
) *SchedulerService {
	return &SchedulerService{
		context:      context,
		config:       config,
		taskRepo:     taskRepo,
		masterClient: masterClient,
	}
}

type SchedulerServiceInterface interface {
	Run()
	PerformingTasks()
	Wait()
}

func (ss *SchedulerService) Run() {
	go ss.PerformingTasks()
}

func (ss *SchedulerService) Wait() {
	ss.wg.Wait()
}

func (ss *SchedulerService) PerformingTasks() {
	ticker := time.NewTicker(60 * ss.config.Scheduler.Duration)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			tasks, err := ss.taskRepo.FindPageByStatus(models.Reporting, 1, 10)
			if err != nil {
				log.Printf("Failed to find tasks: %v", err)
				continue
			}
			if len(tasks) == 0 {
				continue
			}
			ss.wg.Add(len(tasks))
			for _, task := range tasks {
				go func() {
					defer ss.wg.Done()
					req := client.UpdateResultTaskRequest{
						TrackID: task.TrackID,
						Result:  *task.Result,
					}
					statusCode, err := ss.masterClient.UpdateResultTask(ss.config.Master.Url, req)
					if err != nil {
						log.Printf("Failed update result task: %v", err)
						return
					}
					if *statusCode != 200 {
						log.Println("Failed update result task: ", statusCode)
						return
					}
					if _, err := ss.taskRepo.Update(task.TrackID, *task.Result, models.Completed); err != nil {
						log.Printf("Failed to update task: %v", err)
					}
				}()
			}

		case <-ss.context.Done():
			return
		}
	}
}
