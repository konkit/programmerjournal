package database

import (
	"gorm.io/gorm"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
)

type EntryService struct {
	db *gorm.DB
}

func NewEntryService(db *gorm.DB) *EntryService {
	return &EntryService{db}
}

func (es *EntryService) GetEntryByID(entryID uint64) (entry.Entry, error) {
	t := entry.Entry{ID: uint(entryID)}
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

func (es *EntryService) InsertEntry(e *entry.Entry) error {
	var nextRank int64
	es.db.Model(entry.Entry{}).
		Where("created_date = ?", e.CreatedDate).
		Where("rank >= 0").
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
	// TODO: Handle Rank change if the task is deleted in the middle
	err := es.db.Delete(&entry.Entry{}, entryID).Error
	if err != nil {
		return err
	}

	return deleteTagsFromEntry(es.db, entryID)
}

func (es *EntryService) FindTasksFromLastWeek(firstDayOfWeek date.DayDate) ([]entry.Entry, error) {
	monDate := firstDayOfWeek
	tueDate := monDate.PlusDays(1)
	wedDate := tueDate.PlusDays(1)
	thuDate := wedDate.PlusDays(1)
	friDate := thuDate.PlusDays(1)
	satDate := friDate.PlusDays(1)
	sunDate := satDate.PlusDays(1)

	var tasksFromDB []entry.Entry
	err := es.db.Model(entry.Entry{}).
		Where("created_date IN (?, ?, ?, ?, ?, ?, ?)", monDate.Value, tueDate.Value, wedDate.Value, thuDate.Value, friDate.Value, satDate.Value, sunDate.Value).
		Where("status LIKE ?", "Task%").
		Find(&tasksFromDB).
		Error

	if err != nil {
		return nil, err
	}
	return tasksFromDB, nil
}

func (es *EntryService) FindNotesFromLastWeek(firstDayOfWeek date.DayDate) ([]entry.Entry, error) {
	monDate := firstDayOfWeek
	tueDate := monDate.PlusDays(1)
	wedDate := tueDate.PlusDays(1)
	thuDate := wedDate.PlusDays(1)
	friDate := thuDate.PlusDays(1)
	satDate := friDate.PlusDays(1)
	sunDate := satDate.PlusDays(1)

	notesFromDB := []entry.Entry{}
	err := es.db.Model(entry.Entry{}).
		Where("created_date IN (?, ?, ?, ?, ?, ?, ?)", monDate.Value, tueDate.Value, wedDate.Value, thuDate.Value, friDate.Value, satDate.Value, sunDate.Value).
		Where("status LIKE ?", "Note%").
		Find(&notesFromDB).
		Error

	if err != nil {
		return nil, err
	}
	return notesFromDB, nil
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
