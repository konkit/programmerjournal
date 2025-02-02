package database

import (
	"gorm.io/gorm"
	"programmerjournal-backend/model/recurringtask"
)

type RecurringTaskService struct {
	db *gorm.DB
}

func NewRecurringTaskService(db *gorm.DB) *RecurringTaskService {
	return &RecurringTaskService{db: db}
}

func (rs *RecurringTaskService) FindAll() ([]recurringtask.RecurringTask, error) {
	var recurrList []recurringtask.RecurringTask
	err := rs.db.Model(recurringtask.RecurringTask{}).
		Find(&recurrList).
		Error
	return recurrList, err
}

func (rs *RecurringTaskService) Create(newRTask recurringtask.RecurringTask) error {
	return rs.db.Create(&newRTask).Error
}

func (rs *RecurringTaskService) Delete(id uint64) error {
	return rs.db.Delete(&recurringtask.RecurringTask{}, id).Error
}
