package task

import (
	"context"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/entry"

	"github.com/danielgtaylor/huma/v2"
)

type CancelTaskInput struct {
	ID uint `path:"id" example:"123" doc:"ID of the task entry"`
}

type CancelTaskResponse struct {
	Status int
}

func CancelTaskHandler(api huma.API, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "CancelTask",
		Method:      http.MethodPatch,
		Path:        "/api/tasks/{id}/cancelTask",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *CancelTaskInput) (*CancelTaskResponse, error) {
		resp := &CancelTaskResponse{}
		err := CancelTask(es, input.ID)
		if err != nil {
			slog.Error("Error in CancelTaskHandler", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func CancelTask(es *database.EntryService, entryID uint) error {
	t, err := es.GetEntryByID(entryID)
	if err != nil {
		return err
	}

	t.Status = entry.StatusTaskCancelled

	err = es.UpdateEntry(&t)
	if err != nil {
		return err
	}

	err = MoveToTheTop(es, t)
	if err != nil {
		return err
	}

	return nil
}
