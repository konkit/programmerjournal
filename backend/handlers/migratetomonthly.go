package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/date"
	"programmerjournal-backend/taskrepository"
)

type MigrateToMonthlyLogInput struct {
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		Date string `json:"date" doc:"Date when the task should be snoozed"`
	}
}

type MigrateToMonthlyLogResponse struct {
	Status int
}

func MigrateToMonthlyLog(api huma.API, taskRepo *taskrepository.DBRepository) {
	op := huma.Operation{
		OperationID: "MigrateToMonthlyLog",
		Method:      http.MethodPatch,
		Path:        "/api/tasks/{id}/migrateToMonthly",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *MigrateToMonthlyLogInput) (*MigrateToMonthlyLogResponse, error) {
		resp := &MigrateToMonthlyLogResponse{}
		err := taskRepo.MigrateToMonthly(input.ID, date.Parse(input.Body.Date))
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
