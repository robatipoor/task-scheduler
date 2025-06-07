package models

import "gorm.io/gorm"

type Task struct {
	gorm.Model
	TrackUID    string `gorm:"unique;not null"`
	Input       uint   `gorm:"not null"`
	Priority    uint   `gorm:"not null"`
	ScheduleUID *string
	AssignTasks []AssignTask
}

func (Task) TableName() string {
	return "master_tasks"
}
