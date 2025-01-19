package recurringtask

import (
	"gorm.io/gorm"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Create(newRTask RecurringTask) error {
	return s.db.Create(&newRTask).Error
}

func (s *Service) List() ([]RecurringTask, error) {
	var entriesFromDB []RecurringTask
	err := s.db.Model(RecurringTask{}).
		Find(&entriesFromDB).
		Error

	return entriesFromDB, err
}

func (s *Service) Get(newRTask RecurringTask) error {
	return s.db.Create(&newRTask).Error
}

func (s *Service) Update(id uint64, rTask RecurringTask) error {
	rTask.ID = uint(id)
	return s.db.Save(rTask).Error
}

func (s *Service) Delete(id uint64) error {
	return s.db.Delete(&RecurringTask{}, id).Error
}
