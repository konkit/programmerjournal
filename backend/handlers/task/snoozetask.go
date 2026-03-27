package task

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"

	"github.com/danielgtaylor/huma/v2"
)

type SnoozeTaskInput struct {
	ID   uint `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		Date string `json:"date" doc:"Date when the task should be snoozed"`
	}
}

type SnoozeTaskResponse struct {
	Status int
}

func SnoozeTaskHandler(api huma.API, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "SnoozeTask",
		Method:      http.MethodPatch,
		Path:        "/api/tasks/{id}/snooze",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *SnoozeTaskInput) (*SnoozeTaskResponse, error) {
		resp := &SnoozeTaskResponse{}
		dateType := date.GetDateType(input.Body.Date)
		if dateType == date.DateTypeUnrecognized {
			resp.Status = http.StatusBadRequest
			return nil, fmt.Errorf("createdDate in unrecognized date format: %s", input.Body.Date)
		}
		err := SnoozeTask(es, input.ID, date.DateString(input.Body.Date))
		if err != nil {
			slog.Error("Error in SnoozeTaskHandler", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func SnoozeTask(es *database.EntryService, entryID uint, date date.DateString) error {
	snoozedTask, err := es.GetEntryByID(entryID)
	if err != nil {
		return err
	}

	snoozeAfterTaskDate, err := date.IsAfter(snoozedTask.CreatedDate)
	if err != nil {
		return err
	}

	if !snoozeAfterTaskDate {
		return fmt.Errorf("snooze date must be in the future")
	}

	if snoozedTask.Status != entry.StatusTaskCreated {
		return fmt.Errorf("invalid entry status, can only snooze created tasks")
	}

	snoozedTask.Status = entry.StatusTaskSnoozed
	snoozedTask.TaskSnoozedUntil = date
	err = es.UpdateEntry(&snoozedTask)
	if err != nil {
		return err
	}

	err = MoveToTheBottom(es, snoozedTask)
	if err != nil {
		return err
	}

	newTask := entry.Clone(snoozedTask)
	newTask.Status = entry.StatusTaskCreated
	newTask.CreatedDate = date
	return es.InsertEntry(&newTask)
}
