package entryhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/entry"
)

type SetDescriptionInput struct {
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		//Date string `json:"date" doc:"Date when the task should be snoozed"`
		Description string `json:"description" doc:"New title of the task"`
	}
}

type SetTaskDescriptionResponse struct {
	Status int
}

func SetDescription(api huma.API, taskRepo *entry.DBRepository) {
	op := huma.Operation{
		OperationID: "SetDescription",
		Method:      http.MethodPatch,
		Path:        "/api/entries/{id}/setDescription",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *SetDescriptionInput) (*SetTaskDescriptionResponse, error) {
		resp := &SetTaskDescriptionResponse{}
		err := taskRepo.SetTaskDescription(input.ID, input.Body.Description)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
