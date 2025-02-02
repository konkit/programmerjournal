package database

import (
	"gorm.io/gorm"
	"programmerjournal-backend/model/entry"
)

func GetEntryByID(db *gorm.DB, entryID uint64) (entry.Entry, error) {
	t := entry.Entry{ID: uint(entryID)}
	if err := db.First(&t).Error; err != nil {
		return t, err
	}
	return t, nil
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
