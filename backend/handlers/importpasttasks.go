package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/taskrepository"
)

type ImportPastTasksInput struct {
	Date string `path:"date"`
}

type ImportPastTasksResponse struct {
	Status int
}

func ImportPastTasks(api huma.API, taskRepo *taskrepository.DBRepository) {
	op := huma.Operation{
		OperationID: "ImportPastTasks",
		Method:      http.MethodPost,
		Path:        "/api/tasks/importPastTasks/{date}",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *ImportPastTasksInput) (*ImportPastTasksResponse, error) {
		resp := &ImportPastTasksResponse{}
		err := taskRepo.ImportPastTasks(input.Date)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
