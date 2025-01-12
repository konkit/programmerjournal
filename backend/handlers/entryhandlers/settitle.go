package entryhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/entry"
)

type SetTitleInput struct {
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		//Date string `json:"date" doc:"Date when the task should be snoozed"`
		Title string `json:"title" doc:"New title of the task"`
	}
}

type SetTaskTitleResponse struct {
	Status int
}

func SetTitle(api huma.API, taskRepo *entry.Service) {
	op := huma.Operation{
		OperationID: "SetTitle",
		Method:      http.MethodPatch,
		Path:        "/api/entries/{id}/setTitle",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *SetTitleInput) (*SetTaskTitleResponse, error) {
		resp := &SetTaskTitleResponse{}
		err := taskRepo.SetTitle(input.ID, input.Body.Title)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
