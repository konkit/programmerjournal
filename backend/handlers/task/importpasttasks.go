package task

import (
	"context"
	"fmt"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"

	"github.com/danielgtaylor/huma/v2"
)

type ImportPastTasksInput struct {
	Date string `path:"date"`
}

type ImportPastTasksResponse struct {
	Status int
}

func ImportPastTasksHandler(api huma.API, rt *database.RecurringTaskService, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "ImportPastTasks",
		Method:      http.MethodPost,
		Path:        "/api/tasks/pastTasks/{date}/import",
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

			err = importPastTasksFromDay(es, rt, dayDate)
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

			err = importPastTasksFromMonth(es, monthDate)
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

func importPastTasksFromDay(es *database.EntryService, rs *database.RecurringTaskService, today date.DayDate) error {
	dates := getPastDays(today, 30)
	err := importPastCreatedEntries(es, dates, today.Value)
	if err != nil {
		return err
	}

	err = importRecurringTasks(es, rs, today)
	if err != nil {
		return err
	}

	return nil
}

func importPastTasksFromMonth(es *database.EntryService, today date.MonthDate) error {
	dates := getPastMonths(today, 12)
	return importPastCreatedEntries(es, dates, today.Value)
}

func importPastCreatedEntries(es *database.EntryService, dates []date.DateString, targetDate date.DateString) error {
	for _, d := range dates {
		tasks, err := es.FindEntriesByDate(d)
		if err != nil {
			return err
		}

		for _, t := range tasks {
			if t.Status == entry.StatusTaskCreated {
				newTask := entry.Clone(t)
				newTask.CreatedDate = targetDate
				newTask.TaskUpdate = ""
				err := es.InsertEntry(&newTask)
				if err != nil {
					return err
				}

				t.Status = entry.StatusTaskSnoozed
				err = es.UpdateEntry(&t)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func importRecurringTasks(es *database.EntryService, rs *database.RecurringTaskService, today date.DayDate) error {
	recurrList, err := rs.FindAll()
	if err != nil {
		return err
	}

	for _, recurrT := range recurrList {
		if recurrT.DayWithinDate(today) {
			title := fmt.Sprintf("%s - %s", recurrT.TaskTitle, today.Value)

			existingCount, err := es.CountByDateAndRecurringTaskID(today, recurrT.ID)
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
				err := es.InsertEntry(&newTask)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

type PastTasksCountResponse struct {
	Status int
	Body   struct {
		Count int `json:"count"`
	}
}

func CountPastTasks(api huma.API, rt *database.RecurringTaskService, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "CountPastTasks",
		Method:      http.MethodGet,
		Path:        "/api/tasks/pastTasks/{date}/count",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *ImportPastTasksInput) (*PastTasksCountResponse, error) {
		dateType := date.GetDateType(input.Date)

		var count int
		var err error

		switch dateType {
		case date.DateTypeDay:
			dayDate, errParse := date.ParseDayDate(date.DateString(input.Date))
			if errParse != nil {
				return nil, errParse
			}
			count, err = countPendingImportsFromDay(es, rt, dayDate)
		case date.DateTypeMonth:
			monthDate, errParse := date.ParseMonthDate(date.DateString(input.Date))
			if errParse != nil {
				return nil, errParse
			}
			count, err = countPendingImportsFromMonth(es, monthDate)
		case date.DateTypeWeek:
			weekDate, errParse := date.ParseWeekDate(date.DateString(input.Date))
			if errParse != nil {
				return nil, errParse
			}
			count, err = countPendingImportsFromWeek(es, weekDate)
		default:
			return nil, fmt.Errorf("unrecognized date format: %s", input.Date)
		}

		if err != nil {
			return nil, err
		}

		resp := &PastTasksCountResponse{}
		resp.Status = http.StatusOK
		resp.Body.Count = count
		return resp, nil
	})
}

func countPendingImportsFromDay(es *database.EntryService, rs *database.RecurringTaskService, today date.DayDate) (int, error) {
	dates := getPastDays(today, 30)
	count, err := countPastCreatedEntries(es, dates)
	if err != nil {
		return 0, err
	}

	recurrList, err := rs.FindAll()
	if err != nil {
		return 0, err
	}

	for _, recurrT := range recurrList {
		if recurrT.DayWithinDate(today) {
			existingCount, err := es.CountByDateAndRecurringTaskID(today, recurrT.ID)
			if err != nil {
				return 0, err
			}

			if existingCount == 0 {
				count++
			}
		}
	}

	return count, nil
}

func countPendingImportsFromMonth(es *database.EntryService, today date.MonthDate) (int, error) {
	dates := getPastMonths(today, 12)
	return countPastCreatedEntries(es, dates)
}

func countPendingImportsFromWeek(es *database.EntryService, today date.WeekDate) (int, error) {
	dates := getPastWeeks(today, 4)
	return countPastCreatedEntries(es, dates)
}

func countPastCreatedEntries(es *database.EntryService, dates []date.DateString) (int, error) {
	count := 0
	for _, d := range dates {
		tasks, err := es.FindEntriesByDate(d)
		if err != nil {
			return 0, err
		}

		for _, t := range tasks {
			if t.Status == entry.StatusTaskCreated {
				count++
			}
		}
	}
	return count, nil
}

func getPastDays(start date.DayDate, limit int) []date.DateString {
	dates := make([]date.DateString, 0, limit)
	for i := 1; i < limit; i++ {
		dates = append(dates, start.MinusDays(i).Value)
	}
	return dates
}

func getPastWeeks(start date.WeekDate, limit int) []date.DateString {
	dates := make([]date.DateString, 0, limit)
	for i := 1; i < limit; i++ {
		dates = append(dates, start.MinusWeek(i).Value)

	}
	return dates
}

func getPastMonths(start date.MonthDate, limit int) []date.DateString {
	dates := make([]date.DateString, 0, limit)
	for i := 1; i < limit; i++ {
		dates = append(dates, start.MinusMonth(i).Value)
	}
	return dates
}
