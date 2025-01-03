package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/taskrepository"
)

type ChangeTaskRankInput struct {
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		NewRank int `json:"newRank"`
	}
}

type ChangeTaskRankResponse struct {
	Status int
}

func ChangeTaskRank(api huma.API, taskRepo *taskrepository.DBRepository) {
	huma.Patch(api, "/api/tasks/{id}/changeRank", func(ctx context.Context, input *ChangeTaskRankInput) (*ChangeTaskRankResponse, error) {
		resp := &ChangeTaskRankResponse{}
		err := taskRepo.ChangeRank(input.ID, input.Body.NewRank)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
