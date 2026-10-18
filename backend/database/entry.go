package database

import (
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"

	"gorm.io/gorm"
)

type EntryService struct {
	db *gorm.DB
}

func NewEntryService(db *gorm.DB) *EntryService {
	return &EntryService{db}
}

func (es *EntryService) GetEntryByID(entryID uint) (entry.Entry, error) {
	t := entry.Entry{ID: entryID}
	if err := es.db.First(&t).Error; err != nil {
		return t, err
	}
	return t, nil
}

func (es *EntryService) FindEntriesByDate(date date.DateString) ([]entry.Entry, error) {
	var entriesFromDB []entry.Entry
	err := es.db.Model(entry.Entry{}).
		Order("rank").
		Where("created_date = ?", date).
		Find(&entriesFromDB).
		Error

	if err != nil {
		return nil, err
	}

	return entriesFromDB, nil
}
func (es *EntryService) FindActiveEntriesByDate(date date.DateString) ([]entry.Entry, error) {
	var entriesFromDB []entry.Entry
	err := es.db.Model(entry.Entry{}).
		Order("rank").
		Where("created_date = ? AND status NOT IN ('TaskMigrated', 'TaskCancelled', 'TaskSnoozed')", date).
		Find(&entriesFromDB).
		Error

	if err != nil {
		return nil, err
	}

	return entriesFromDB, nil
}
func (es *EntryService) InsertEntry(e *entry.Entry) error {
	var nextRank int64
	es.db.Model(entry.Entry{}).
		Where("created_date = ?", e.CreatedDate).
		Count(&nextRank)
	e.Rank = int(nextRank)

	err := es.db.Create(e).Error
	if err != nil {
		return err
	}

	err = saveTags(es.db, e)
	if err != nil {
		return err
	}

	return nil
}

func (es *EntryService) UpdateEntry(e *entry.Entry) error {
	err := es.db.Save(e).Error
	if err != nil {
		return err
	}

	err = saveTags(es.db, e)
	if err != nil {
		return err
	}

	return nil
}

func (es *EntryService) DeleteEntry(entryID uint) error {
	err := es.db.Delete(&entry.Entry{}, entryID).Error
	if err != nil {
		return err
	}

	return deleteTagsFromEntry(es.db, entryID)
}

func (es *EntryService) CountByDateAndRecurringTaskID(date date.DayDate, recurringTaskID uint) (int64, error) {
	var existingCount int64
	err := es.db.Model(entry.Entry{}).
		Where("created_date = ?", date.Value).
		Where("recurring_task_id = ?", recurringTaskID).
		Count(&existingCount).
		Error
	return existingCount, err
}

func (es *EntryService) FindByDateAndTaskID(date date.DateString, taskID string) (*entry.Entry, error) {
	t := entry.Entry{}
	err := es.db.Model(entry.Entry{}).
		Where("created_date = ?", date).
		Where("task_id = ?", taskID).
		First(&t).
		Error
	return &t, err
}

func (es *EntryService) FindTasksByTaskID(taskID string) ([]entry.Entry, error) {
	var tasksFromDB []entry.Entry
	err := es.db.Model(entry.Entry{}).
		Where("task_id = ?", taskID).
		Find(&tasksFromDB).
		Error

	return tasksFromDB, err
}

func (es *EntryService) FindAllNotes() ([]entry.Entry, error) {
	var notesFromDB []entry.Entry
	err := es.db.Model(entry.Entry{}).
		Where("status = ?", entry.StatusNote).
		Order("created_date desc").
		Find(&notesFromDB).
		Error

	if err != nil {
		return nil, err
	}
	return notesFromDB, nil
}
