package repositories

import (
	"database/sql"
	"fmt"

	"github.com/robatipoor/task-scheduler/internal/master/models"
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
	FindByTrackUID(trackUID string) (*models.Task, error)
	Take(page, pageSize int, scheduleUID string) ([]models.Task, error)
	Save(trackID string, input, priority uint) (*models.Task, error)
}

func (tr *TaskRepository) FindByID(id uint) (*models.Task, error) {
	var task models.Task
	err := tr.db.First(&task, id).Error
	return &task, err
}

func (tr *TaskRepository) FindByTrackUID(trackUID string) (*models.Task, error) {
	var task models.Task
	err := tr.db.Where("track_uid = ?", trackUID).Preload("AssignTasks").First(&task).Error
	return &task, err
}

func (tr *TaskRepository) Take(page, pageSize int, scheduleUID string) ([]models.Task, error) {
	var tasks []models.Task
	offset := (page - 1) * pageSize
	err := tr.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ").Error; err != nil {
			return err
		}
		if err := tx.Where("schedule_uid IS NULL").Order("priority DESC").Offset(offset).Limit(pageSize).Find(&tasks).Error; err != nil {
			return err
		}
		if len(tasks) == 0 {
			return nil
		}
		for i := range tasks {
			tasks[i].ScheduleUID = &scheduleUID
		}
		if err := tx.Save(&tasks).Error; err != nil {
			return fmt.Errorf("error updating tasks: %w", err)
		}
		return nil
	}, &sql.TxOptions{Isolation: sql.LevelRepeatableRead})

	return tasks, err
}

func (tr *TaskRepository) Save(trackID string, input, priority uint) (*models.Task, error) {
	task := models.Task{
		TrackUID:    trackID,
		Input:       input,
		Priority:    priority,
		ScheduleUID: nil,
	}
	err := tr.db.Create(&task).Error
	return &task, err
}
