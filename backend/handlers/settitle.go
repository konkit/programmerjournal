package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/taskrepository"
)

type SetTaskTitleInput struct {
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		//Date string `json:"date" doc:"Date when the task should be snoozed"`
		Title string `json:"title" doc:"New title of the task"`
	}
}

type SetTaskTitleResponse struct {
	Status int
}

func SetTaskTitle(api huma.API, taskRepo *taskrepository.DBRepository) {
	op := huma.Operation{
		OperationID: "SetTaskTitle",
		Method:      http.MethodPatch,
		Path:        "/api/tasks/{id}/setTitle",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *SetTaskTitleInput) (*SetTaskTitleResponse, error) {
		resp := &SetTaskTitleResponse{}
		err := taskRepo.SetTaskTitle(input.ID, input.Body.Title)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
