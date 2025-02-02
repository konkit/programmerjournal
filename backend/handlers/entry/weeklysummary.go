package entry

import (
	"context"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
)

type WeeklyTaskSummaryInput struct {
	Date string `path:"date" example:"2024-05-05" doc:"First day of the week to summarize"`
}

type WeeklyTaskSummaryOutput struct {
	Body entry.WeeklySummary
}

func WeeklySummaryHandler(api huma.API, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "WeeklySummary",
		Method:      http.MethodGet,
		Path:        "/api/entries/weeklySummary/{date}",
		Tags:        []string{"EntryService"},
	}
	huma.Register(api, op, func(ctx context.Context, input *WeeklyTaskSummaryInput) (*WeeklyTaskSummaryOutput, error) {
		resp := &WeeklyTaskSummaryOutput{}

		dayDate, err := date.ParseDayDate(date.DateString(input.Date))
		if err != nil {
			return nil, err
		}
		summ, err := WeeklySummary(es, dayDate)
		if err != nil {
			return nil, err
		}
		resp.Body = summ
		return resp, nil
	})
}

func WeeklySummary(es *database.EntryService, firstDayOfWeek date.DayDate) (entry.WeeklySummary, error) {
	isDateMonday := checkIfDateIsMonday(firstDayOfWeek)
	if !isDateMonday {
		return entry.WeeklySummary{}, fmt.Errorf("the selected date is not the first day of the week")
	}

	tasksFromDB, err := es.FindTasksFromLastWeek(firstDayOfWeek)
	if err != nil {
		return entry.WeeklySummary{}, fmt.Errorf("failed to get weekly summary: %v", err)
	}

	taskMap := groupByTaskID(tasksFromDB)

	summaryArr := []entry.TaskSummary{}
	for _, taskArr := range taskMap {
		t, err := findLastDayTask(taskArr)
		if err != nil {
			return entry.WeeklySummary{}, err
		}
		updates := getUpdates(taskArr)

		summary := entry.TaskSummary{
			TaskEntry: t,
			Updates:   updates,
		}
		summaryArr = append(summaryArr, summary)
	}

	notesFromDB, err := es.FindNotesFromLastWeek(firstDayOfWeek)
	if err != nil {
		return entry.WeeklySummary{}, fmt.Errorf("failed to get weekly summary: %v", err)
	}

	ws := entry.WeeklySummary{
		TaskSummaries: summaryArr,
		Notes:         notesFromDB,
	}

	return ws, err
}

func checkIfDateIsMonday(d date.DayDate) bool {
	// TODO: Check if monday
	return true
}

// Function to group a slice of structs by a specific property
func groupByTaskID(tasks []entry.Entry) map[string][]entry.Entry {
	grouped := make(map[string][]entry.Entry)
	for _, tt := range tasks {
		grouped[tt.TaskID] = append(grouped[tt.TaskID], tt)
	}
	return grouped
}

func findLastDayTask(arr []entry.Entry) (entry.Entry, error) {
	lastTask := entry.Entry{
		CreatedDate: "1000-01-01",
	}
	for _, t := range arr {
		createdDate, err := date.ParseDayDate(t.CreatedDate)
		if err != nil {
			return entry.Entry{}, fmt.Errorf("entries dates not in day format: %s", t.CreatedDate)
		}
		lastTaskDate, err := date.ParseDayDate(lastTask.CreatedDate)
		if err != nil {
			return entry.Entry{}, fmt.Errorf("entries dates not in day format: %s", lastTask.CreatedDate)
		}
		isAfter, err := createdDate.IsAfter(lastTaskDate)
		if err != nil {
			return entry.Entry{}, err
		}

		if isAfter {
			lastTask = t
		}
	}

	return lastTask, nil
}

func getUpdates(arr []entry.Entry) []entry.TaskUpdate {
	var res []entry.TaskUpdate
	for _, t := range arr {
		tu := entry.TaskUpdate{
			Date:   t.CreatedDate,
			Update: t.TaskUpdate,
			Status: t.Status,
		}
		res = append(res, tu)
	}
	return res
}
