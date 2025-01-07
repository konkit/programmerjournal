package taskhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/entry"
)

type GetTaskSummaryInput struct {
	ID uint64 `path:"id" example:"123" doc:"ID of the task entry"`
}

type GetTaskSummaryResponse struct {
	Status int
	Body   *entry.TaskSummary
}

func GetTaskSummary(api huma.API, taskRepo *entry.DBRepository) {
	op := huma.Operation{
		OperationID: "GetTaskSummary",
		Method:      http.MethodGet,
		Path:        "/api/tasks/{id}/summary",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *GetTaskSummaryInput) (*GetTaskSummaryResponse, error) {
		resp := &GetTaskSummaryResponse{}
		summary, err := taskRepo.GetTaskSummary(input.ID)
		if err != nil {
			return nil, err
		}

		resp.Body = summary
		resp.Status = http.StatusOK
		return resp, nil
	})
}
