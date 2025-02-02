package task

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"net/http"
	"programmerjournal-backend/database"
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

func MigrateTaskToMonthlyLogHandler(api huma.API, db *gorm.DB) {
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

		err = MigrateToMonthly(db, input.ID, monthDate)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func MigrateToMonthly(db *gorm.DB, entryID uint64, date date.MonthDate) error {
	t, err := database.GetEntryByID(db, entryID)
	if err != nil {
		return err
	}

	t.Status = entry.StatusTaskMigrated
	t.TaskSnoozedUntil = date.Value
	err = database.UpdateEntry(db, &t)
	if err != nil {
		return err
	}

	ee, err := findByDateAndTaskID(db, date.Value, t.TaskID)
	if err != nil {
		return err
	}
	if ee != nil {
		ee.Status = entry.StatusTaskCreated
		//db.Save(&ee)
		err = database.UpdateEntry(db, ee)
		if err != nil {
			return err
		}
	} else {
		//nextRank := utils.FetchNextRank(db, date.Value)

		newTask := entry.Clone(t)
		newTask.Status = entry.StatusTaskCreated
		newTask.CreatedDate = date.Value
		//newTask.Rank = nextRank

		err = database.InsertEntry(db, &newTask)
		if err != nil {
			return err
		}
	}

	return nil
}
