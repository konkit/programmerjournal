package taskhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
)

type MigrateTaskToDailyLogInput struct {
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		Date string `json:"date" doc:"Date when the task should be snoozed"`
	}
}

type MigrateTaskToDailyLogResponse struct {
	Status int
}

func MigrateTaskToDailyLog(api huma.API, taskRepo *entry.Service) {
	op := huma.Operation{
		OperationID: "MigrateTaskToDailyLog",
		Method:      http.MethodPatch,
		Path:        "/api/tasks/{id}/migrateToDaily",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *MigrateTaskToDailyLogInput) (*MigrateTaskToDailyLogResponse, error) {
		resp := &MigrateTaskToDailyLogResponse{}

		dayDate, err := date.ParseDayDate(date.DateString(input.Body.Date))
		if err != nil {
			return nil, err
		}

		err = taskRepo.MigrateToDaily(input.ID, dayDate)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
