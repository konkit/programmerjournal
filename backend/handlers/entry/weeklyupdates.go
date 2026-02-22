package entry

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"

	"github.com/danielgtaylor/huma/v2"
)

type WeeklyTaskUpdatesInput struct {
	Date string `path:"date" example:"2024-05-05" doc:"First day of the week to summarize"`
}

type WeeklyTaskUpdatesOutput struct {
	Body map[date.DateString][]entry.Entry
}

func WeeklyUpdatesHandler(api huma.API, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "WeeklyUpdates",
		Method:      http.MethodGet,
		Path:        "/api/entries/weeklyUpdates/{date}",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *WeeklyTaskUpdatesInput) (*WeeklyTaskUpdatesOutput, error) {
		resp := &WeeklyTaskUpdatesOutput{}

		dayDate, err := date.ParseDayDate(date.DateString(input.Date))
		if err != nil {
			return nil, err
		}
		summ, err := FetchWeeklyUpdates(es, dayDate)
		if err != nil {
			slog.Error("Error in WeeklyUpdatesHandler", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
		}
		resp.Body = summ
		return resp, nil
	})
}

func FetchWeeklyUpdates(es *database.EntryService, firstDayOfWeek date.DayDate) (map[date.DateString][]entry.Entry, error) {
	isDateMonday := checkIfDateIsMonday(firstDayOfWeek)
	if !isDateMonday {
		return nil, fmt.Errorf("the selected date is not the first day of the week")
	}

	tasksFromDB, err := es.FindTasksFromLastWeek(firstDayOfWeek)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly summary: %v", err)
	}

	filteredTasks := []entry.Entry{}
	for _, task := range tasksFromDB {
		if task.Status == entry.StatusTaskCreated && task.TaskUpdate == "" {
			continue
		}

		if task.Status == entry.StatusTaskMigrated && task.TaskUpdate == "" {
			continue
		}

		if task.Status == entry.StatusTaskSnoozed && task.TaskUpdate == "" {
			continue
		}

		filteredTasks = append(filteredTasks, task)
	}

	dateToTaskUpdate := make(map[date.DateString][]entry.Entry)
	for _, task := range filteredTasks {
		dateToTaskUpdate[task.CreatedDate] = append(dateToTaskUpdate[task.CreatedDate], task)
	}

	return dateToTaskUpdate, nil
}
