package database

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"programmerjournal-backend/model/entry"
	"programmerjournal-backend/model/recurringtask"
)

func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&entry.Entry{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate table for Entry")
	}
	err = db.AutoMigrate(&recurringtask.RecurringTask{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate table for RecurringTask")
	}
	return db, nil
}
