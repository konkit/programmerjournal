package utils

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

func FetchNextRank(db *gorm.DB, date date.DateString) int {
	var count int64
	db.Model(entry.Entry{}).
		Where("created_date = ?", date).
		Where("rank >= 0").
		Count(&count)
	return int(count)
}
