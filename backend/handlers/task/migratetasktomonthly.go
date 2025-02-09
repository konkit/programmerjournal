package task

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
)

type MigrateTaskToMonthlyLogInput struct {
	ID   uint `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		Date string `json:"date" doc:"Date when the task should be snoozed"`
	}
}

type MigrateTaskToMonthlyLogResponse struct {
	Status int
}

func MigrateTaskToMonthlyLogHandler(api huma.API, es *database.EntryService) {
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

		err = MigrateToMonthly(es, input.ID, monthDate)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func MigrateToMonthly(es *database.EntryService, entryID uint, date date.MonthDate) error {
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
