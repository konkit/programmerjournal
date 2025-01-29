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

type ImportPastTasksInput struct {
	Date string `path:"date"`
}

type ImportPastTasksResponse struct {
	Status int
}

func ImportPastTasksHandler(api huma.API, db *gorm.DB) {
	op := huma.Operation{
		OperationID: "ImportPastTasks",
		Method:      http.MethodPost,
		Path:        "/api/tasks/importPastTasks/{date}",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *ImportPastTasksInput) (*ImportPastTasksResponse, error) {
		resp := &ImportPastTasksResponse{}

		dateType := date.GetDateType(input.Date)

		switch dateType {
		case date.DateTypeDay:
			dayDate, err := date.ParseDayDate(date.DateString(input.Date))
			if err != nil {
				resp.Status = http.StatusBadRequest
				return nil, err
			}

			err = ImportPastTasksFromDay(db, dayDate)
			if err != nil {
				return nil, err
			}

			resp.Status = http.StatusOK
			return resp, nil
		case date.DateTypeMonth:
			monthDate, err := date.ParseMonthDate(date.DateString(input.Date))
			if err != nil {
				resp.Status = http.StatusBadRequest
				return nil, err
			}

			err = ImportPastTasksFromMonth(db, monthDate)
			if err != nil {
				return nil, err
			}

			resp.Status = http.StatusOK
			return resp, nil
		default:
			resp.Status = http.StatusBadRequest
			return nil, fmt.Errorf("unrecognized date format: %s", input.Date)
		}
	})
}

func ImportPastTasksFromDay(db *gorm.DB, today date.DayDate) error {
	for i := 1; i < 30; i++ {
		current := today.MinusDays(i)

		tasks, err := ListDayEntries(db, current)
		if err != nil {
			return err
		}

		for _, t := range tasks {
			if t.Status == entry.StatusTaskCreated {
				nextRank := fetchNextRank(db, today.Value)

				newTask := entry.Clone(t)
				newTask.CreatedDate = today.Value
				newTask.Rank = nextRank
				newTask.TaskUpdate = ""
				err := db.Save(&newTask).Error
				if err != nil {
					return err
				}

				t.Status = entry.StatusTaskSnoozed
				err = db.Save(&t).Error
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func ImportPastTasksFromMonth(db *gorm.DB, today date.MonthDate) error {
	for i := 1; i < 12; i++ {
		current := today.MinusMonth(i)

		tasks, err := ListMonthEntries(db, current)
		if err != nil {
			return err
		}

		for _, t := range tasks {
			if t.Status == entry.StatusTaskCreated {
				nextRank := fetchNextRank(db, today.Value)

				newTask := entry.Clone(t)
				newTask.CreatedDate = today.Value
				newTask.Rank = nextRank
				newTask.TaskUpdate = ""
				db.Save(&newTask)

				t.Status = entry.StatusTaskSnoozed
				db.Save(&t)
			}
		}
	}

	return nil
}

func ListDayEntries(db *gorm.DB, date date.DayDate) ([]entry.Entry, error) {
	return listEntries(db, date.Value)
}

func ListMonthEntries(db *gorm.DB, date date.MonthDate) ([]entry.Entry, error) {
	return listEntries(db, date.Value)
}

func listEntries(db *gorm.DB, date date.DateString) ([]entry.Entry, error) {
	var entriesFromDB []entry.Entry
	err := db.Model(entry.Entry{}).
		Order("rank").
		Where("created_date = ?", date).
		Find(&entriesFromDB).
		Error

	return entriesFromDB, err
}
