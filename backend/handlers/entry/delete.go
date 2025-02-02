package entry

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/database"
)

type DeleteEntryInput struct {
	ID uint `path:"id" example:"123" doc:"ID of the task entry"`
}

type DeleteEntryResponse struct {
	Status int
}

func DeleteEntryHandler(api huma.API, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "DeleteEntry",
		Method:      http.MethodDelete,
		Path:        "/api/entries/{id}/delete",
		Tags:        []string{"EntryService"},
	}
	huma.Register(api, op, func(ctx context.Context, input *DeleteEntryInput) (*DeleteEntryResponse, error) {
		resp := &DeleteEntryResponse{}
		err := es.DeleteEntry(input.ID)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
