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
	if t.Rank <= 0 {
		// Do not move priority entries.
		return nil
	}
	entriesFromDB, err := es.FindEntriesByDate(t.CreatedDate)
	if err != nil {
		return err
	}

	firstNotDoneIndex, err := getFirstNotDoneIndex(entriesFromDB, t)
	if err != nil {
		return err
	}

	err = entryhandlers.ChangeRank(es, t.ID, firstNotDoneIndex)
	if err != nil {
		return err
	}

	return nil
}

func MoveToTheBottom(es *database.EntryService, t entry.Entry) error {
	if t.Rank < 0 {
		// Do not move priority entries.
		return nil
	}
	entriesFromDB, err := es.FindEntriesByDate(t.CreatedDate)
	if err != nil {
		return err
	}

	lastIndex := len(entriesFromDB) - 1

	err = entryhandlers.ChangeRank(es, t.ID, lastIndex)
	if err != nil {
		return err
	}

	return nil
}

func getFirstNotDoneIndex(entriesFromDB []entry.Entry, current entry.Entry) (int, error) {
	i := 0
	for i = 0; i < len(entriesFromDB); i++ {
		e := entriesFromDB[i]

		if e.ID == current.ID {
			return i, nil
		}

		if !(e.Status == entry.StatusTaskSnoozed || e.Status == entry.StatusTaskMigrated || e.Status == entry.StatusTaskDone) {
			return i, nil
		}
	}

	return i, nil
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
