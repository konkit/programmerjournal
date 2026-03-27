package task

import (
	"context"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/entry"

	"github.com/danielgtaylor/huma/v2"
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
			slog.Error("Error in SetTaskDoneHandler", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
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
		err = MoveToTheBottom(es, t)
		if err != nil {
			return err
		}
	}

	return nil
}
