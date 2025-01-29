package task

import (
	"context"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
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

func SnoozeTaskHandler(api huma.API, db *gorm.DB) {
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
		err := SnoozeTask(db, input.ID, date.DateString(input.Body.Date))
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func SnoozeTask(db *gorm.DB, entryID uint64, date date.DateString) error {
	snoozedTask, err := getEntryByID(db, entryID)
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
	db.Save(&snoozedTask)

	nextRank := fetchNextRank(db, date)

	newTask := entry.Clone(snoozedTask)
	newTask.Status = entry.StatusTaskCreated
	newTask.CreatedDate = date
	newTask.Rank = nextRank
	db.Save(&newTask)

	return nil
}
