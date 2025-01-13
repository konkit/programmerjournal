package recurringtask

import (
	"gorm.io/gorm"
)

type RecurringTaskService struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *RecurringTaskService {
	return &RecurringTaskService{db: db}
}

func (s *RecurringTaskService) Create(newRTask RecurringTask) error {
	return s.db.Create(&newRTask).Error
}
