package models

import (
	"gorm.io/gorm"
)

type AssingStatus uint

const (
	Assinged AssingStatus = iota
	Submitted
	Completed
	Failed
)

type AssignTask struct {
	gorm.Model
	TaskID       uint         `gorm:"not null"`
	WorkerID     uint         `gorm:"not null"`
	WorkerUrl    string       `gorm:"not null"`
	Status       AssingStatus `gorm:"not null"`
	ErrorMessage string
	Result       *uint
}
