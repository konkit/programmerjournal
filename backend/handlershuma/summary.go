package handlershuma

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
)

type GetTaskSummaryInput struct {
	ID uint64 `path:"id" example:"123" doc:"ID of the task entry"`
}

type GetTaskSummaryResponse struct {
	Status int
	Body   *task.TaskSummary
}

func GetTaskSummary(api huma.API, taskRepo *taskrepository.DBRepository) {
	huma.Get(api, "/api/tasks/{id}/summary", func(ctx context.Context, input *GetTaskSummaryInput) (*GetTaskSummaryResponse, error) {
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
