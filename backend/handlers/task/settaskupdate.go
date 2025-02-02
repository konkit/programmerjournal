package task

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/database"
)

type SetTaskUpdateInput struct {
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		//Date string `json:"date" doc:"Date when the task should be snoozed"`
		Update string `json:"update" doc:"New title of the task"`
	}
}

type SetTaskUpdateResponse struct {
	Status int
}

func SetTaskUpdateHandler(api huma.API, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "SetTaskUpdate",
		Method:      http.MethodPatch,
		Path:        "/api/tasks/{id}/setUpdate",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *SetTaskUpdateInput) (*SetTaskUpdateResponse, error) {
		resp := &SetTaskUpdateResponse{}
		err := SetTaskUpdate(es, input.ID, input.Body.Update)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func SetTaskUpdate(es *database.EntryService, entryID uint64, update string) error {
	t, err := es.GetEntryByID(entryID)
	if err != nil {
		return err
	}

	t.TaskUpdate = update

	return es.UpdateEntry(&t)
}
