package taskhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
)

type ImportPastTasksInput struct {
	Date string `path:"date"`
}

type ImportPastTasksResponse struct {
	Status int
}

func ImportPastTasks(api huma.API, taskRepo *entry.DBRepository) {
	op := huma.Operation{
		OperationID: "ImportPastTasks",
		Method:      http.MethodPost,
		Path:        "/api/tasks/importPastTasks/{date}",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *ImportPastTasksInput) (*ImportPastTasksResponse, error) {
		resp := &ImportPastTasksResponse{}
		err := taskRepo.ImportPastTasks(date.Parse(input.Date))
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
