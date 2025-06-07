package repositories

import (
	"github.com/robatipoor/task-scheduler/internal/worker/models"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

type TaskRepositoryInterface interface {
	FindByID(id uint) (*models.Task, error)
	FindByTrackID(TrackID uint) (*models.Task, error)
	FindPageByStatus(status models.TaskStatus, page, pageSize int) ([]models.Task, error)
	Save(TrackID, input uint) (*models.Task, error)
	Update(TrackID, result uint, status models.TaskStatus) (*models.Task, error)
}

func (tr *TaskRepository) FindByID(id uint) (*models.Task, error) {
	var task models.Task
	err := tr.db.First(&task, id).Error
	return &task, err
}

func (tr *TaskRepository) FindByTrackID(TrackID uint) (*models.Task, error) {
	var task models.Task
	err := tr.db.Where("track_id = ?", TrackID).First(&task).Error
	return &task, err
}

func (tr *TaskRepository) FindPageByStatus(status models.TaskStatus, page, pageSize int) ([]models.Task, error) {
	var tasks []models.Task
	offset := (page - 1) * pageSize
	err := tr.db.Where("status = ?", status).Offset(offset).Limit(pageSize).Find(&tasks).Error
	return tasks, err
}

func (tr *TaskRepository) Save(TrackID, input uint) (*models.Task, error) {
	task := models.Task{
		TrackID: TrackID,
		Input:   input,
	}
	err := tr.db.Create(&task).Error
	return &task, err
}

func (tr *TaskRepository) Update(TrackID, result uint, status models.TaskStatus) (*models.Task, error) {
	task, err := tr.FindByTrackID(TrackID)
	if err != nil {
		return nil, err
	}
	task.Status = status
	task.Result = &result
	if err := tr.db.Save(&task).Error; err != nil {
		return nil, err
	}
	return task, nil
}
