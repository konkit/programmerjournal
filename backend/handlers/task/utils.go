package task

import (
	"errors"
	"programmerjournal-backend/database"
	entryhandlers "programmerjournal-backend/handlers/entry"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"

	"gorm.io/gorm"
)

func MoveToTheTop(es *database.EntryService, t entry.Entry) error {
	return entryhandlers.ChangeRank(es, t.ID, 0)
}

func MoveToTheBottom(es *database.EntryService, t entry.Entry) error {
	entriesFromDB, err := es.FindEntriesByDate(t.CreatedDate)
	if err != nil {
		return err
	}

	lastIndex := len(entriesFromDB) - 1
	return entryhandlers.ChangeRank(es, t.ID, lastIndex)
}

func findByDateAndTaskID(es *database.EntryService, date date.DateString, taskID string) (*entry.Entry, error) {
	t, err := es.FindByDateAndTaskID(date, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return t, nil
}
