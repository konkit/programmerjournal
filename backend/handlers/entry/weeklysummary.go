package entry

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
	"sort"

	"github.com/danielgtaylor/huma/v2"
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
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *WeeklyTaskSummaryInput) (*WeeklyTaskSummaryOutput, error) {
		resp := &WeeklyTaskSummaryOutput{}

		dayDate, err := date.ParseDayDate(date.DateString(input.Date))
		if err != nil {
			return nil, err
		}
		summ, err := WeeklySummary(es, dayDate)
		if err != nil {
			slog.Error("Error in WeeklySummaryHandler", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
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

	finishedThisWeekMap := make(map[date.DateString][]entry.TaskSummary)
	var updatedThisWeek []entry.TaskSummary
	var others []entry.TaskSummary

	for _, summ := range summaryArr {
		finished := false
		var finishedDate date.DateString
		for _, upd := range summ.Updates {
			if upd.Status == entry.StatusTaskDone {
				finishedDate = upd.Date
				finished = true
			}
		}
		if finished {
			finishedThisWeekMap[finishedDate] = append(finishedThisWeekMap[finishedDate], summ)
			continue
		}

		updated := false
		for _, upd := range summ.Updates {
			if len(upd.Update) > 0 {
				updated = true
			}
		}
		if updated {
			updatedThisWeek = append(updatedThisWeek, summ)
			continue
		}

		others = append(others, summ)
	}

	var finishedThisWeek []entry.FinishedThisWeekSummary
	var mapKeys []string
	for k := range finishedThisWeekMap {
		mapKeys = append(mapKeys, string(k))
	}
	//sort.Slice(mapKeys)
	sort.Sort(sort.Reverse(sort.StringSlice(mapKeys)))
	for _, d := range mapKeys {
		dateString := date.DateString(d)
		f := entry.FinishedThisWeekSummary{
			Date:    dateString,
			Updates: finishedThisWeekMap[dateString],
		}
		finishedThisWeek = append(finishedThisWeek, f)
	}

	notesFromDB, err := es.FindNotesFromLastWeek(firstDayOfWeek)
	if err != nil {
		return entry.WeeklySummary{}, fmt.Errorf("failed to get weekly summary: %v", err)
	}

	ws := entry.WeeklySummary{
		TaskSummaries:         summaryArr,
		TasksFinishedThisWeek: finishedThisWeek,
		TasksUpdatedThisWeek:  updatedThisWeek,
		OtherTasks:            others,
		Notes:                 notesFromDB,
	}

	return ws, err
}

func filterFinishedThisWeek(input []entry.TaskSummary) []entry.TaskSummary {
	var result []entry.TaskSummary
	for _, summ := range input {
		for _, upd := range summ.Updates {
			if upd.Status == entry.StatusTaskDone {
				result = append(result, summ)
			}
		}
	}
	return result
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
