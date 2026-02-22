package task

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"

	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
)

type MigrateTaskToDailyLogInput struct {
	ID   uint `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		Date string `json:"date" doc:"Date when the task should be snoozed"`
	}
}

type MigrateTaskToDailyLogResponse struct {
	Status int
}

func MigrateTaskToDailyLogHandler(api huma.API, es *database.EntryService) {
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

		err = MigrateToDaily(es, input.ID, dayDate)
		if err != nil {
			slog.Error("Error in MigrateTaskToDailyLogHandler", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func MigrateToDaily(es *database.EntryService, entryID uint, date date.DayDate) error {
	t, err := es.GetEntryByID(entryID)
	if err != nil {
		return err
	}

	t.Status = entry.StatusTaskSnoozed // Change the name to "migrate lower" or "migrated to more specific" something like this
	t.TaskSnoozedUntil = date.Value
	err = es.UpdateEntry(&t)
	if err != nil {
		return err
	}

	err = MoveToTheTop(es, t)
	if err != nil {
		return err
	}

	ee, err := findByDateAndTaskID(es, date.Value, t.TaskID)
	if err != nil {
		return err
	}
	if ee != nil {
		ee.Status = entry.StatusTaskCreated
		err = es.UpdateEntry(ee)
		if err != nil {
			return err
		}
	} else {
		newTask := entry.Clone(t)
		newTask.Status = entry.StatusTaskCreated
		newTask.CreatedDate = date.Value
		err = es.InsertEntry(&newTask)
		if err != nil {
			return err
		}
	}

	return nil
}

func findByDateAndTaskID(es *database.EntryService, date date.DateString, taskID string) (*entry.Entry, error) {
	t, err := es.FindByDateAndTaskID(date, taskID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return t, nil
}
