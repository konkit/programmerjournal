package handlershuma

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/taskrepository"
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

func SetTaskUpdate(api huma.API, taskRepo *taskrepository.DBRepository) {
	huma.Patch(api, "/api/tasks/{id}/setUpdate", func(ctx context.Context, input *SetTaskUpdateInput) (*SetTaskUpdateResponse, error) {
		resp := &SetTaskUpdateResponse{}
		err := taskRepo.SetTaskUpdate(input.ID, input.Body.Update)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
