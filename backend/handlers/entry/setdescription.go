package entry

import (
	"context"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"

	"github.com/danielgtaylor/huma/v2"
)

type SetDescriptionInput struct {
	ID   uint `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		//Date string `json:"date" doc:"Date when the task should be snoozed"`
		Description string `json:"description" doc:"New title of the task"`
	}
}

type SetTaskDescriptionResponse struct {
	Status int
}

func SetDescriptionHandler(api huma.API, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "SetDescription",
		Method:      http.MethodPatch,
		Path:        "/api/entries/{id}/setDescription",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *SetDescriptionInput) (*SetTaskDescriptionResponse, error) {
		resp := &SetTaskDescriptionResponse{}
		err := SetDescription(es, input.ID, input.Body.Description)
		if err != nil {
			slog.Error("Error in SetDescriptionHandler", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func SetDescription(es *database.EntryService, entryID uint, description string) error {
	t, err := es.GetEntryByID(entryID)
	if err != nil {
		return err
	}

	t.Description = description

	return es.UpdateEntry(&t)
}
