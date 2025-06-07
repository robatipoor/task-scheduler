package repositories

import (
	"errors"

	"github.com/robatipoor/task-scheduler/internal/master/models"
	"gorm.io/gorm"
)

type WorkerRepository struct {
	db *gorm.DB
}

func NewWorkerRepository(db *gorm.DB) *WorkerRepository {
	return &WorkerRepository{db: db}
}

type WorkerRepositoryInterface interface {
	FindByUrl(url string) (*models.Worker, error)
	FindActiveWorkers() ([]models.Worker, error)
	SaveOrUpdate(url string, status models.WorkerStatus) (*models.Worker, error)
}

func (wr *WorkerRepository) FindByUrl(url string) (*models.Worker, error) {
	var worker models.Worker
	err := wr.db.Where("url = ?", url).First(&worker).Error
	return &worker, err
}

func (wr *WorkerRepository) FindActiveWorkers() ([]models.Worker, error) {
	var worker []models.Worker
	err := wr.db.Where("status = ?", models.Active).Find(&worker).Error
	return worker, err
}

func (wr *WorkerRepository) SaveOrUpdate(url string, status models.WorkerStatus) (*models.Worker, error) {
	worker, err := wr.FindByUrl(url)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			worker = &models.Worker{
				Url:    url,
				Status: status,
			}
			if err := wr.db.Create(worker).Error; err != nil {
				return nil, err
			}
			return worker, nil
		}
		return nil, err
	}
	worker.Status = status
	if err := wr.db.Save(&worker).Error; err != nil {
		return nil, err
	}
	return worker, nil
}
