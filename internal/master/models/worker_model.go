package models

import "gorm.io/gorm"

type WorkerStatus uint

const (
	Active WorkerStatus = iota
	InActive
)

type Worker struct {
	gorm.Model
	Url         string       `gorm:"not null"`
	Status      WorkerStatus `gorm:"not null"`
	AssignTasks []AssignTask
}
