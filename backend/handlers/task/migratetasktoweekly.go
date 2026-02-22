package task

import (
	"context"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"

	"github.com/danielgtaylor/huma/v2"
)

type MigrateTaskToWeeklyLogInput struct {
	ID   uint `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		Date string `json:"date" doc:"Date when the task should be snoozed"`
	}
}

type MigrateTaskToWeeklyLogResponse struct {
	Status int
}

func MigrateTaskToWeeklyLogHandler(api huma.API, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "MigrateTaskToWeeklyLog",
		Method:      http.MethodPatch,
		Path:        "/api/tasks/{id}/migrateToWeekly",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *MigrateTaskToWeeklyLogInput) (*MigrateTaskToWeeklyLogResponse, error) {
		resp := &MigrateTaskToWeeklyLogResponse{}

		weekDate, err := date.ParseWeekDate(date.DateString(input.Body.Date))
		if err != nil {
			return nil, err
		}

		err = MigrateToWeekly(es, input.ID, weekDate)
		if err != nil {
			slog.Error("Error in MigrateTaskToWeeklyLogHandler", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func MigrateToWeekly(es *database.EntryService, entryID uint, date date.WeekDate) error {
	t, err := es.GetEntryByID(entryID)
	if err != nil {
		return err
	}

	t.Status = entry.StatusTaskMigrated
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
