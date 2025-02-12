package task

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/database"
	entryhandlers "programmerjournal-backend/handlers/entry"
	"programmerjournal-backend/model/entry"
)

type SetTaskDoneInput struct {
	ID   uint `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		//Date string `json:"date" doc:"Date when the task should be snoozed"`
		Done bool `json:"done" doc:"If task should be set as done or not"`
	}
}

type SetTaskDoneResponse struct {
	Status int
}

func SetTaskDoneHandler(api huma.API, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "SetTaskDone",
		Method:      http.MethodPatch,
		Path:        "/api/tasks/{id}/setDone",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *SetTaskDoneInput) (*SetTaskDoneResponse, error) {
		resp := &SetTaskDoneResponse{}
		err := SetTaskDone(es, input.ID, input.Body.Done)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func SetTaskDone(es *database.EntryService, entryID uint, done bool) error {
	t, err := es.GetEntryByID(entryID)
	if err != nil {
		return err
	}

	if done == true {
		t.Status = entry.StatusTaskDone
	} else {
		t.Status = entry.StatusTaskCreated
	}

	err = es.UpdateEntry(&t)
	if err != nil {
		return err
	}

	if done == true {
		err = MoveToTheTop(es, t)
		if err != nil {
			return err
		}
	}

	return nil
}

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
