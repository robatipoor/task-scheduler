package services

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/robatipoor/task-scheduler/internal/master/client"
	"github.com/robatipoor/task-scheduler/internal/master/config"
	"github.com/robatipoor/task-scheduler/internal/master/models"
	"github.com/robatipoor/task-scheduler/internal/master/repositories"
)

type SchedulerService struct {
	context        context.Context
	wg             sync.WaitGroup
	config         *config.Configure
	workerClient   client.WorkerClientInterface
	workerRepo     repositories.WorkerRepositoryInterface
	assignTaskRepo repositories.AssingTaskRepositoryInterface
	taskRepo       repositories.TaskRepositoryInterface
}

func NewSchedulerService(
	context context.Context,
	config *config.Configure,
	workerClient client.WorkerClientInterface,
	workerRepo repositories.WorkerRepositoryInterface,
	taskRepo repositories.TaskRepositoryInterface,
	assignTaskRepo repositories.AssingTaskRepositoryInterface,
) SchedulerService {
	return SchedulerService{
		context:        context,
		config:         config,
		workerClient:   workerClient,
		workerRepo:     workerRepo,
		taskRepo:       taskRepo,
		assignTaskRepo: assignTaskRepo,
	}
}

type SchedulerServiceInterface interface {
	Run()
	HeartbeatCheck()
	AssignTasks(scheduleUID string)
	DistributeTasks()
	RecoverTasks()
	Wait()
}

func (ss *SchedulerService) Run() {
	go ss.HeartbeatCheck()
	go ss.AssignTasks(uuid.New().String())
	go ss.DistributeTasks()
	go ss.RecoverTasks()
}

func (ss *SchedulerService) Wait() {
	ss.wg.Wait()
}

func (ss *SchedulerService) RecoverTasks() {
	ss.wg.Add(1)
	d := 1800 * ss.config.Scheduler.Duration
	ticker := time.NewTicker(d)
	defer func() {
		ticker.Stop()
		ss.wg.Done()
	}()
	for {
		select {
		case <-ticker.C:
			tasks, err := ss.assignTaskRepo.FindAllUnknownState(time.Now().Add(-1 * d))
			if err != nil {
				log.Printf("get unknown state tasks failed: %v\n", err)
				continue
			}
			for _, atask := range tasks {
				result, err := ss.workerClient.GetResultTask(atask.WorkerUrl, atask.ID)
				var status models.AssingStatus
				var errorMessage string
				if err != nil {
					log.Printf("assinged task to the worker failed: %v\n", err)
					status = models.Failed
					errorMessage = err.Error()
				} else if *result.Status == 2 {
					status = models.Completed
				} else if *result.Status == 3 {
					status = models.Failed
				} else {
					continue
				}

				_, err = ss.assignTaskRepo.Update(atask.ID, status, result.Result, errorMessage)

				if err != nil {
					log.Printf("update status submitted task failed: %v\n", err)
					continue
				}
			}
		case <-ss.context.Done():
			return
		}
	}
}

func (ss *SchedulerService) DistributeTasks() {
	ss.wg.Add(1)
	ticker := time.NewTicker(10 * ss.config.Scheduler.Duration)
	defer func() {
		ticker.Stop()
		ss.wg.Done()
	}()
	for {
		select {
		case <-ticker.C:
			tasks, err := ss.assignTaskRepo.FindPageByStatus(models.Submitted, 1, 100)
			if err != nil {
				log.Printf("get submitted tasks failed: %v\n", err)
				continue
			}
			if len(tasks) == 0 {
				log.Println("no submitted task exists")
				continue
			}
			for _, atask := range tasks {
				task, err := ss.taskRepo.FindByID(atask.TaskID)
				if err != nil {
					log.Printf("get task: %d failed: %v\n", atask.ID, err)
					continue
				}
				req := client.AssingTaskWorkerRequest{
					TrackID: atask.ID,
					Input:   task.Input,
				}
				statusCode, err := ss.workerClient.AssignTask(atask.WorkerUrl, req)
				var status models.AssingStatus
				var errorMessage string
				if err != nil {
					log.Printf("assinged task to the worker failed: %v\n", err)
					status = models.Failed
					errorMessage = err.Error()
				} else if !(*statusCode == 201 || *statusCode == 200) {
					status = models.Failed
				} else {
					status = models.Assinged
				}

				_, err = ss.assignTaskRepo.Update(req.TrackID, status, nil, errorMessage)

				if err != nil {
					log.Printf("update status submitted task failed: %v\n", err)
					continue
				}
			}
		case <-ss.context.Done():
			return
		}
	}
}

func (ss *SchedulerService) AssignTasks(scheduleUID string) {
	ss.wg.Add(1)
	ticker := time.NewTicker(10 * ss.config.Scheduler.Duration)
	defer func() {
		ticker.Stop()
		ss.wg.Done()
	}()

	for {
		select {
		case <-ticker.C:
			workers, err := ss.workerRepo.FindActiveWorkers()

			if err != nil {
				log.Printf("find active workers failed: %v\n", err)
				continue
			}

			workersLen := len(workers)
			if workersLen == 0 {
				log.Printf("no active worker found\n")
				continue
			}

			wgen := func() func(task models.Task) (uint, string) {
				i := 0
				return func(task models.Task) (uint, string) {
					if i < workersLen-1 {
						i++
					} else {
						i = 0
					}
					return workers[i].ID, workers[i].Url
				}
			}

			_, err = ss.taskRepo.Assign(1, 100, scheduleUID, wgen)
			if err != nil {
				log.Printf("take tasks failed: %v\n", err)
				continue
			}

		case <-ss.context.Done():
			return
		}
	}
}

func (ss *SchedulerService) HeartbeatCheck() {
	ss.wg.Add(1)
	ticker := time.NewTicker(30 * ss.config.Scheduler.Duration)
	defer func() {
		ticker.Stop()
		ss.wg.Done()
	}()
	for {
		select {
		case <-ticker.C:
			if err := ss.checkActiveWorkers(); err != nil {
				log.Printf("Heartbeat check failed: %v\n", err)
			}
		case <-ss.context.Done():
			return
		}
	}
}

func (ss *SchedulerService) checkActiveWorkers() error {
	workers, err := ss.workerRepo.FindActiveWorkers()
	if err != nil {
		return fmt.Errorf("failed to get active workers: %w", err)
	}

	for _, worker := range workers {
		if err := ss.checkWorkerStatus(worker); err != nil {
			log.Printf("Error checking status for worker %s: %v\n", worker.Url, err)
		}
	}
	return nil
}

func (ss *SchedulerService) checkWorkerStatus(worker models.Worker) error {
	status, err := ss.workerClient.HealthCheck(worker.Url)
	if err != nil || *status != 200 {
		log.Printf("Worker %s seems to be inactive. Health check status: %v, error: %v\n", worker.Url, status, err)
		_, err := ss.workerRepo.SaveOrUpdate(worker.Url, models.InActive)
		return err
	} else {
		_, err := ss.workerRepo.SaveOrUpdate(worker.Url, models.Active)
		return err
	}
}
