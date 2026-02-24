package task

import (
	"programmerjournal-backend/database"
	entryhandlers "programmerjournal-backend/handlers/entry"
	"programmerjournal-backend/model/entry"
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
