package entryhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/entry"
)

type DeleteEntryInput struct {
	ID uint64 `path:"id" example:"123" doc:"ID of the task entry"`
}

type DeleteEntryResponse struct {
	Status int
}

func DeleteEntry(api huma.API, taskRepo *entry.Service) {
	op := huma.Operation{
		OperationID: "DeleteEntry",
		Method:      http.MethodDelete,
		Path:        "/api/entries/{id}/delete",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *DeleteEntryInput) (*DeleteEntryResponse, error) {
		resp := &DeleteEntryResponse{}
		err := taskRepo.Delete(input.ID)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
