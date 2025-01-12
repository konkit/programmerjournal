package entryhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/entry"
)

type UpdateEntryInput struct {
	Body entry.Entry
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
}

type UpdateEntryResponse struct {
	Status int
}

func UpdateEntry(api huma.API, taskRepo *entry.Service) {
	op := huma.Operation{
		OperationID: "UpdateEntry",
		Method:      http.MethodPut,
		Path:        "/api/tasks/{id}/update",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *UpdateEntryInput) (*UpdateEntryResponse, error) {
		resp := &UpdateEntryResponse{}
		err := taskRepo.Update(input.ID, input.Body)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
