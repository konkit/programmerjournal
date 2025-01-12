package taskhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
)

type MigrateTaskToMonthlyLogInput struct {
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		Date string `json:"date" doc:"Date when the task should be snoozed"`
	}
}

type MigrateTaskToMonthlyLogResponse struct {
	Status int
}

func MigrateTaskToMonthlyLog(api huma.API, taskRepo *entry.Service) {
	op := huma.Operation{
		OperationID: "MigrateTaskToMonthlyLog",
		Method:      http.MethodPatch,
		Path:        "/api/tasks/{id}/migrateToMonthly",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *MigrateTaskToMonthlyLogInput) (*MigrateTaskToMonthlyLogResponse, error) {
		resp := &MigrateTaskToMonthlyLogResponse{}

		monthDate, err := date.ParseMonthDate(date.DateString(input.Body.Date))
		if err != nil {
			return nil, err
		}

		err = taskRepo.MigrateToMonthly(input.ID, monthDate)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
