package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/taskrepository"
)

type SetTaskDescriptionInput struct {
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		//Date string `json:"date" doc:"Date when the task should be snoozed"`
		Description string `json:"description" doc:"New title of the task"`
	}
}

type SetTaskDescriptionResponse struct {
	Status int
}

func SetTaskDescription(api huma.API, taskRepo *taskrepository.DBRepository) {
	op := huma.Operation{
		OperationID: "SetTaskDescription",
		Method:      http.MethodPatch,
		Path:        "/api/tasks/{id}/setDescription",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *SetTaskDescriptionInput) (*SetTaskDescriptionResponse, error) {
		resp := &SetTaskDescriptionResponse{}
		err := taskRepo.SetTaskDescription(input.ID, input.Body.Description)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
