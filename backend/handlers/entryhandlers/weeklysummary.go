package entryhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
)

type WeeklyTaskSummaryInput struct {
	Date string `path:"date" example:"2024-05-05" doc:"First day of the week to summarize"`
}

type WeeklyTaskSummaryOutput struct {
	Body entry.WeeklySummary
}

func WeeklySummary(api huma.API, taskRepo *entry.Service) {
	op := huma.Operation{
		OperationID: "WeeklySummary",
		Method:      http.MethodGet,
		Path:        "/api/entries/weeklySummary/{date}",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *WeeklyTaskSummaryInput) (*WeeklyTaskSummaryOutput, error) {
		resp := &WeeklyTaskSummaryOutput{}

		dayDate, err := date.ParseDayDate(date.DateString(input.Date))
		if err != nil {
			return nil, err
		}
		summ, err := taskRepo.WeeklySummary(dayDate)
		if err != nil {
			return nil, err
		}
		resp.Body = summ
		return resp, nil
	})
}
