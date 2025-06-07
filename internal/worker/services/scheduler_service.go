package services

import (
	"context"
	"log"
	"math/rand"
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
) SchedulerService {
	return SchedulerService{
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
			tasks, err := ss.taskRepo.FindPageByStatus(models.Submitted, 1, 10)
			if err != nil {
				log.Printf("Failed to find tasks: %v", err)
				continue
			}
			for _, task := range tasks {
				ss.wg.Add(1)
				go func() {
					defer ss.wg.Done()
					time.Sleep(time.Duration(task.Input) * time.Millisecond)
					r := uint(rand.Intn(int(task.Input)))
					_, err := ss.taskRepo.Update(task.TrackID, r, models.Completed)
					if err != nil {
						log.Printf("Failed update task %v", err)
						return
					}
					req := client.UpdateResultTaskRequest{
						TrackID: task.TrackID,
						Result:  r,
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
				}()
			}

		case <-ss.context.Done():
			return
		}
	}
}
