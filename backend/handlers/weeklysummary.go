package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/date"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
)

type WeeklySummaryOutput struct {
	Body []task.TaskWeeklySummary
}

func WeeklySummary(api huma.API, taskRepo *taskrepository.DBRepository) {
	op := huma.Operation{
		OperationID: "WeeklySummary",
		Method:      http.MethodGet,
		Path:        "/api/tasks/weeklySummary/{date}",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *struct {
		Date string `path:"date" example:"2024-05-05" doc:"First day of the week to summarize"`
	}) (*WeeklySummaryOutput, error) {
		resp := &WeeklySummaryOutput{}
		summ, err := taskRepo.WeeklySummary(date.Parse(input.Date))
		if err != nil {
			return nil, err
		}
		resp.Body = summ
		return resp, nil
	})
}
