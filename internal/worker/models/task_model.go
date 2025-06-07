package models

import "gorm.io/gorm"

type TaskStatus uint

const (
	Submitted = iota
	Running
	Completed
	Failed
)

type Task struct {
	gorm.Model
	TrackID uint `gorm:"unique;not null"`
	Input   uint `gorm:"not null"`
	Result  *uint
	Status  TaskStatus `gorm:"not null"`
}

func (Task) TableName() string {
	return "worker_tasks"
}
