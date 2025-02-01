package task

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
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

func MigrateTaskToDailyLogHandler(api huma.API, db *gorm.DB) {
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

		err = MigrateToDaily(db, input.ID, dayDate)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func MigrateToDaily(db *gorm.DB, entryID uint64, date date.DayDate) error {
	t, err := getEntryByID(db, entryID)
	if err != nil {
		return err
	}

	t.Status = entry.StatusTaskMigrated
	t.TaskSnoozedUntil = date.Value
	db.Save(&t)

	ee, err := findByDateAndTaskID(db, date.Value, t.TaskID)
	if err != nil {
		return err
	}
	if ee != nil {
		ee.Status = entry.StatusTaskCreated
		db.Save(&ee)
	} else {
		nextRank := fetchNextRank(db, date.Value)

		newTask := entry.Clone(t)
		newTask.Status = entry.StatusTaskCreated
		newTask.CreatedDate = date.Value
		newTask.Rank = nextRank
		db.Save(&newTask)
	}

	return nil
}

func findByDateAndTaskID(db *gorm.DB, date date.DateString, taskID string) (*entry.Entry, error) {
	t := entry.Entry{}
	err := db.Model(entry.Entry{}).
		Where("created_date = ?", date).
		Where("task_id = ?", taskID).
		First(&t).
		Error

	if err != nil {
		if err.Error() == "record not found" {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &t, nil
}

func getEntryByID(db *gorm.DB, entryID uint64) (entry.Entry, error) {
	t := entry.Entry{ID: uint(entryID)}
	if err := db.First(&t).Error; err != nil {
		return t, err
	}
	return t, nil
}

//func fetchNextRank(db *gorm.DB, date date.DateString) int {
//	var count int64
//	db.Model(entry.Entry{}).
//		Where("created_date = ?", date).
//		Where("rank >= 0").
//		Count(&count)
//	return int(count)
//}
