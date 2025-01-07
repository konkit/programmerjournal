package taskhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
)

type WeeklyTaskSummaryOutput struct {
	Body []entry.TaskSummary
}

func WeeklyTaskSummary(api huma.API, taskRepo *entry.DBRepository) {
	op := huma.Operation{
		OperationID: "WeeklyTaskSummary",
		Method:      http.MethodGet,
		Path:        "/api/tasks/weeklySummary/{date}",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *struct {
		Date string `path:"date" example:"2024-05-05" doc:"First day of the week to summarize"`
	}) (*WeeklyTaskSummaryOutput, error) {
		resp := &WeeklyTaskSummaryOutput{}
		summ, err := taskRepo.WeeklyTaskSummary(date.Parse(input.Date))
		if err != nil {
			return nil, err
		}
		resp.Body = summ
		return resp, nil
	})
}
