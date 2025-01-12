package taskhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
)

type SnoozeTaskInput struct {
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		Date string `json:"date" doc:"Date when the task should be snoozed"`
	}
}

type SnoozeTaskResponse struct {
	Status int
}

func SnoozeTask(api huma.API, taskRepo *entry.Service) {
	op := huma.Operation{
		OperationID: "SnoozeTask",
		Method:      http.MethodPatch,
		Path:        "/api/tasks/{id}/snooze",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *SnoozeTaskInput) (*SnoozeTaskResponse, error) {
		resp := &SnoozeTaskResponse{}
		dateString, err := date.ParseDateString(input.Body.Date)
		if err != nil {
			return nil, err
		}
		err = taskRepo.SnoozeTask(input.ID, dateString)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
