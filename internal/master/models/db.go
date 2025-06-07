package models

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Task{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&AssignTask{})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&Worker{})
	if err != nil {
		return nil, err
	}

	return db, err
}
