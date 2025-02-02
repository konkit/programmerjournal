package task

import (
	"context"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
	"programmerjournal-backend/model/recurringtask"
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
				newTask := entry.Clone(t)
				newTask.CreatedDate = today.Value
				newTask.TaskUpdate = ""
				err := database.InsertEntry(db, &newTask)
				if err != nil {
					return err
				}

				t.Status = entry.StatusTaskSnoozed
				err = database.UpdateEntry(db, &t)
				if err != nil {
					return err
				}
			}
		}
	}

	var recurrList []recurringtask.RecurringTask
	err := db.Model(recurringtask.RecurringTask{}).
		Find(&recurrList).
		Error
	if err != nil {
		return err
	}

	for _, recurrT := range recurrList {
		if recurrT.DayWithinDate(today) {
			title := fmt.Sprintf("%s - %s", recurrT.TaskTitle, today.Value)

			existingCount, err := database.CountByDateAndRecurringTaskID(db, today, recurrT.ID)
			if err != nil {
				return err
			}

			if existingCount == 0 {
				newTask := entry.Entry{
					Title:            title,
					Status:           entry.StatusTaskCreated,
					CreatedDate:      today.Value,
					Description:      recurrT.TaskDescription,
					TaskID:           "",
					TaskUpdate:       "",
					TaskSnoozedUntil: "",
					RecurringTaskID:  recurrT.ID,
				}
				err := database.InsertEntry(db, &newTask)
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
				newTask := entry.Clone(t)
				newTask.CreatedDate = today.Value
				newTask.TaskUpdate = ""
				err = database.InsertEntry(db, &newTask)
				if err != nil {
					return err
				}

				t.Status = entry.StatusTaskSnoozed
				err = database.UpdateEntry(db, &t)
				if err != nil {
					return err
				}

			}
		}
	}

	return nil
}

func ListDayEntries(db *gorm.DB, date date.DayDate) ([]entry.Entry, error) {
	return database.FindEntriesByDate(db, date.Value)
}

func ListMonthEntries(db *gorm.DB, date date.MonthDate) ([]entry.Entry, error) {
	return database.FindEntriesByDate(db, date.Value)
}
