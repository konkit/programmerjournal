package database

import (
	"gorm.io/gorm"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
)

func GetEntryByID(db *gorm.DB, entryID uint64) (entry.Entry, error) {
	t := entry.Entry{ID: uint(entryID)}
	if err := db.First(&t).Error; err != nil {
		return t, err
	}
	return t, nil
}

func FindEntriesByDate(db *gorm.DB, date date.DateString) ([]entry.Entry, error) {
	var entriesFromDB []entry.Entry
	err := db.Model(entry.Entry{}).
		Order("rank").
		Where("created_date = ?", date).
		Find(&entriesFromDB).
		Error

	if err != nil {
		return nil, err
	}

	return entriesFromDB, nil
}

func InsertEntry(db *gorm.DB, e *entry.Entry) error {
	var nextRank int64
	db.Model(entry.Entry{}).
		Where("created_date = ?", e.CreatedDate).
		Where("rank >= 0").
		Count(&nextRank)
	e.Rank = int(nextRank)

	err := db.Create(e).Error
	if err != nil {
		return err
	}

	err = saveTags(db, e)
	if err != nil {
		return err
	}

	return nil
}

func UpdateEntry(db *gorm.DB, e *entry.Entry) error {
	err := db.Save(e).Error
	if err != nil {
		return err
	}

	err = saveTags(db, e)
	if err != nil {
		return err
	}

	return nil
}

func DeleteEntry(db *gorm.DB, entryID uint) error {
	// TODO: Handle Rank change if the task is deleted in the middle
	err := db.Delete(&entry.Entry{}, entryID).Error
	if err != nil {
		return err
	}

	return deleteTagsFromEntry(db, entryID)
}

func FindTasksFromLastWeek(db *gorm.DB, firstDayOfWeek date.DayDate) ([]entry.Entry, error) {
	monDate := firstDayOfWeek
	tueDate := monDate.PlusDays(1)
	wedDate := tueDate.PlusDays(1)
	thuDate := wedDate.PlusDays(1)
	friDate := thuDate.PlusDays(1)
	satDate := friDate.PlusDays(1)
	sunDate := satDate.PlusDays(1)

	var tasksFromDB []entry.Entry
	err := db.Model(entry.Entry{}).
		Where("created_date IN (?, ?, ?, ?, ?, ?, ?)", monDate.Value, tueDate.Value, wedDate.Value, thuDate.Value, friDate.Value, satDate.Value, sunDate.Value).
		Where("status LIKE ?", "Task%").
		Find(&tasksFromDB).
		Error

	if err != nil {
		return nil, err
	}
	return tasksFromDB, nil
}

func FindNotesFromLastWeek(db *gorm.DB, firstDayOfWeek date.DayDate) ([]entry.Entry, error) {
	monDate := firstDayOfWeek
	tueDate := monDate.PlusDays(1)
	wedDate := tueDate.PlusDays(1)
	thuDate := wedDate.PlusDays(1)
	friDate := thuDate.PlusDays(1)
	satDate := friDate.PlusDays(1)
	sunDate := satDate.PlusDays(1)

	notesFromDB := []entry.Entry{}
	err := db.Model(entry.Entry{}).
		Where("created_date IN (?, ?, ?, ?, ?, ?, ?)", monDate.Value, tueDate.Value, wedDate.Value, thuDate.Value, friDate.Value, satDate.Value, sunDate.Value).
		Where("status LIKE ?", "Note%").
		Find(&notesFromDB).
		Error

	if err != nil {
		return nil, err
	}
	return notesFromDB, nil
}

func CountByDateAndRecurringTaskID(db *gorm.DB, date date.DayDate, recurringTaskID uint) (int64, error) {
	var existingCount int64
	err := db.Model(entry.Entry{}).
		Where("created_date = ?", date.Value).
		Where("recurring_task_id = ?", recurringTaskID).
		Count(&existingCount).
		Error
	return existingCount, err
}

func FindByDateAndTaskID(db *gorm.DB, date date.DateString, taskID string) (*entry.Entry, error) {
	t := entry.Entry{}
	err := db.Model(entry.Entry{}).
		Where("created_date = ?", date).
		Where("task_id = ?", taskID).
		First(&t).
		Error
	return &t, err
}

func FindTasksByTaskID(db *gorm.DB, taskID string) ([]entry.Entry, error) {
	var tasksFromDB []entry.Entry
	err := db.Model(entry.Entry{}).
		Where("task_id = ?", taskID).
		Find(&tasksFromDB).
		Error

	return tasksFromDB, err
}
