package handlershuma

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/taskrepository"
)

type SetTaskDoneInput struct {
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		//Date string `json:"date" doc:"Date when the task should be snoozed"`
		Done bool `json:"done" doc:"If task should be set as done or not"`
	}
}

type SetTaskDoneResponse struct {
	Status int
}

func SetTaskDone(api huma.API, taskRepo *taskrepository.DBRepository) {
	huma.Patch(api, "/api/tasks/{id}/setDone", func(ctx context.Context, input *SetTaskDoneInput) (*SetTaskDoneResponse, error) {
		resp := &SetTaskDoneResponse{}
		err := taskRepo.SetTaskDone(input.ID, input.Body.Done)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
