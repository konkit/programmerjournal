package task

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
	"programmerjournal-backend/model/recurringtask"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
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
				slog.Error("Error in ImportPastTasksHandler (day)", "error", err)
				return nil, huma.Error500InternalServerError(err.Error())
			}

			err = importPastTasksFromWeek(es, dayDate.ToWeekDate())
			if err != nil {
				slog.Error("Error in ImportPastTasksHandler (week)", "error", err)
				return nil, huma.Error500InternalServerError(err.Error())
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
				slog.Error("Error in ImportPastTasksHandler", "error", err)
				return nil, huma.Error500InternalServerError(err.Error())
			}

			resp.Status = http.StatusOK
			return resp, nil
		case date.DateTypeWeek:
			weekDate, err := date.ParseWeekDate(date.DateString(input.Date))
			if err != nil {
				resp.Status = http.StatusBadRequest
				slog.Warn("Error in ImportPastTasksHandler", "error", err)
				return nil, err
			}
			err = importPastTasksFromWeek(es, weekDate)
			if err != nil {
				slog.Error("Error in ImportPastTasksHandler", "error", err)
				return nil, huma.Error500InternalServerError(err.Error())
			}

			resp.Status = http.StatusOK
			return resp, nil
		default:
			resp.Status = http.StatusBadRequest
			slog.Warn("unrecognized data format", "input.Date", input.Date)
			return nil, fmt.Errorf("unrecognized date format: %s", input.Date)
		}
	})
}

func importPastTasksFromDay(es *database.EntryService, rs *database.RecurringTaskService, today date.DayDate) error {
	err := unmigrateWeeklyTasks(es, today)
	if err != nil {
		return err
	}

	// Recurring tasks missed on past days are created on their own day first,
	// so that importPastCreatedEntries moves them to today together with other unfinished tasks
	err = importRecurringTasks(es, rs, today)
	if err != nil {
		return err
	}

	dates := getPastDaysPreviousWeekToCurrentWeek(today)
	return importPastCreatedEntries(es, dates, today.Value)
}

func importPastTasksFromMonth(es *database.EntryService, today date.MonthDate) error {
	dates := getPastMonths(today, 12)
	return importPastCreatedEntries(es, dates, today.Value)
}

func importPastTasksFromWeek(es *database.EntryService, thisWeek date.WeekDate) error {
	dates := getPastWeeks(thisWeek, 12)
	return importPastCreatedEntries(es, dates, thisWeek.Value)
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

func unmigrateWeeklyTasks(es *database.EntryService, today date.DayDate) error {
	currentWeek := today.ToWeekDate()
	weeklyTasks, err := es.FindEntriesByDate(currentWeek.Value)
	if err != nil {
		return err
	}

	for _, wt := range weeklyTasks {
		if wt.Status == entry.StatusTaskSnoozed && date.GetDateType(string(wt.TaskSnoozedUntil)) == date.DateTypeDay {
			snoozedUntil, _ := date.ParseDayDate(wt.TaskSnoozedUntil)
			isBefore, _ := snoozedUntil.IsAfter(today)
			if !isBefore && snoozedUntil.Value != today.Value {
				// Task was migrated to a past day.
				// Check if there is an active (Created) daily entry for this TaskID.
				allDailyTasks, err := es.FindTasksByTaskID(wt.TaskID)
				if err != nil {
					continue
				}

				hasActiveDaily := false
				for _, dt := range allDailyTasks {
					if date.GetDateType(string(dt.CreatedDate)) == date.DateTypeDay && dt.Status == entry.StatusTaskCreated {
						hasActiveDaily = true
						// Snooze the daily entry
						dt.Status = entry.StatusTaskSnoozed
						es.UpdateEntry(&dt)
					}
				}

				if hasActiveDaily {
					// Mark weekly task back as created
					wt.Status = entry.StatusTaskCreated
					err = es.UpdateEntry(&wt)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

type recurringTaskInstance struct {
	task recurringtask.RecurringTask
	day  date.DayDate
}

// findMissingRecurringTasks returns the instances of recurring tasks that were due between the start
// of the previous week and today (the same window as for past daily tasks), but were never generated,
// e.g. because the app wasn't opened on that day.
func findMissingRecurringTasks(es *database.EntryService, recurrList []recurringtask.RecurringTask, today date.DayDate) ([]recurringTaskInstance, error) {
	days := append(getPastDaysPreviousWeekToCurrentWeek(today), today.Value)

	var missing []recurringTaskInstance
	for _, recurrT := range recurrList {
		for _, d := range days {
			if recurrT.GeneratedUntil != "" && d <= recurrT.GeneratedUntil {
				continue
			}

			day, err := date.ParseDayDate(d)
			if err != nil {
				return nil, err
			}
			if !recurrT.DayWithinDate(day) {
				continue
			}

			existingCount, err := es.CountByDateAndRecurringTaskID(day, recurrT.ID)
			if err != nil {
				return nil, err
			}
			if existingCount == 0 {
				missing = append(missing, recurringTaskInstance{task: recurrT, day: day})
			}
		}
	}
	return missing, nil
}

func importRecurringTasks(es *database.EntryService, rs *database.RecurringTaskService, today date.DayDate) error {
	recurrList, err := rs.FindAll()
	if err != nil {
		return err
	}

	missing, err := findMissingRecurringTasks(es, recurrList, today)
	if err != nil {
		return err
	}

	for _, instance := range missing {
		newTask := entry.Entry{
			Title:            fmt.Sprintf("%s - %s", instance.task.TaskTitle, instance.day.Value),
			Status:           entry.StatusTaskCreated,
			CreatedDate:      instance.day.Value,
			Description:      instance.task.TaskDescription,
			TaskID:           uuid.NewString(),
			TaskUpdate:       "",
			TaskSnoozedUntil: "",
			RecurringTaskID:  instance.task.ID,
		}
		err := es.InsertEntry(&newTask)
		if err != nil {
			return err
		}
	}

	for _, recurrT := range recurrList {
		if today.Value > recurrT.GeneratedUntil {
			err := rs.SetGeneratedUntil(recurrT.ID, today.Value)
			if err != nil {
				return err
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
			countDay, errDay := countPendingImportsFromDay(es, rt, dayDate)
			if errDay != nil {
				return nil, huma.Error500InternalServerError(errDay.Error())
			}
			countWeek, errWeek := countPendingImportsFromWeek(es, dayDate.ToWeekDate())
			if errWeek != nil {
				return nil, huma.Error500InternalServerError(errWeek.Error())
			}
			count = countDay + countWeek
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
			slog.Error("Error in CountPastTasks", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp := &PastTasksCountResponse{}
		resp.Status = http.StatusOK
		resp.Body.Count = count
		return resp, nil
	})
}

func countPendingImportsFromDay(es *database.EntryService, rs *database.RecurringTaskService, today date.DayDate) (int, error) {
	dates := getPastDaysPreviousWeekToCurrentWeek(today)
	count, err := countPastCreatedEntries(es, dates)
	if err != nil {
		return 0, err
	}

	recurrList, err := rs.FindAll()
	if err != nil {
		return 0, err
	}

	missing, err := findMissingRecurringTasks(es, recurrList, today)
	if err != nil {
		return 0, err
	}

	return count + len(missing), nil
}

func countPendingImportsFromMonth(es *database.EntryService, today date.MonthDate) (int, error) {
	dates := getPastMonths(today, 12)
	return countPastCreatedEntries(es, dates)
}

func countPendingImportsFromWeek(es *database.EntryService, today date.WeekDate) (int, error) {
	dates := getPastWeeks(today, 12)
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

func getPastDaysPreviousWeekToCurrentWeek(today date.DayDate) []date.DateString {
	thisWeek := today.ToWeekDate()
	prevWeek := thisWeek.MinusWeek(1)
	startOfPrevWeek := prevWeek.GetStartDay()

	dates := []date.DateString{}
	current := startOfPrevWeek
	for {
		if current.Value == today.Value {
			break
		}
		dates = append(dates, current.Value)
		current = current.PlusDays(1)
	}

	// Reverse the order to match the behavior of getPastDays (newest first)
	for i, j := 0, len(dates)-1; i < j; i, j = i+1, j-1 {
		dates[i], dates[j] = dates[j], dates[i]
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
