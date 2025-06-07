package repositories

import (
	"time"

	"github.com/robatipoor/task-scheduler/internal/master/models"
	"gorm.io/gorm"
)

type AssingTaskRepository struct {
	db *gorm.DB
}

func NewAssingTaskRepository(db *gorm.DB) *AssingTaskRepository {
	return &AssingTaskRepository{db: db}
}

type AssingTaskRepositoryInterface interface {
	FindByID(id uint) (*models.AssignTask, error)
	FindAllUnknownState(time time.Time) ([]models.AssignTask, error)
	FindPageByStatus(status models.AssingStatus, page, pageSize int) ([]models.AssignTask, error)
	Save(taskID, workerID uint, workerUrl string) (*models.AssignTask, error)
	Update(id uint, status models.AssingStatus, result *uint, errMsg string) (*models.AssignTask, error)
}

func (ar *AssingTaskRepository) FindAllUnknownState(time time.Time) ([]models.AssignTask, error) {
	var tasks []models.AssignTask
	err := ar.db.Where("result IS NULL AND status = ? AND updated_at <= ?", models.Submitted, time).Find(&tasks).Error
	return tasks, err
}

func (ar *AssingTaskRepository) FindPageByStatus(status models.AssingStatus, page, pageSize int) ([]models.AssignTask, error) {
	var tasks []models.AssignTask
	offset := (page - 1) * pageSize
	err := ar.db.Where("status = ?", status).Offset(offset).Limit(pageSize).Find(&tasks).Error
	return tasks, err
}

func (ar *AssingTaskRepository) FindByID(id uint) (*models.AssignTask, error) {
	var task models.AssignTask
	err := ar.db.First(&task, id).Error
	return &task, err
}

func (atr *AssingTaskRepository) Save(taskID, workerID uint, workerUrl string) (*models.AssignTask, error) {
	assignTask := models.AssignTask{
		TaskID:       taskID,
		WorkerID:     workerID,
		WorkerUrl:    workerUrl,
		Status:       models.Submitted,
		ErrorMessage: "",
		Result:       nil,
	}
	err := atr.db.Create(&assignTask).Error
	return &assignTask, err
}

func (atr *AssingTaskRepository) Update(id uint, status models.AssingStatus, result *uint, errMsg string) (*models.AssignTask, error) {
	task, err := atr.FindByID(id)
	if err != nil {
		return nil, err
	}
	task.Status = status
	task.Result = result
	task.ErrorMessage = errMsg
	if err := atr.db.Save(&task).Error; err != nil {
		return nil, err
	}
	return task, nil
}
