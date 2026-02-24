package task

import (
	"fmt"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
	"programmerjournal-backend/model/recurringtask"
	"strconv"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestImportPastTasks(t *testing.T) {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	es := database.NewEntryService(db)
	rts := database.NewRecurringTaskService(db)

	_, api := humatest.New(t)
	ImportPastTasksHandler(api, rts, es)

	createEntry := func(id int, dateParam string, rank int, status entry.Status) entry.Entry {
		return entry.Entry{
			TaskID:      strconv.Itoa(id),
			Title:       fmt.Sprintf("test %d", id),
			Status:      status,
			CreatedDate: date.DateString(dateParam),
			TaskUpdate:  "",
			Rank:        rank,
		}
	}

	testCases := []struct {
		name         string
		today        string
		initTasks    []entry.Entry
		wantResponse []entry.Entry
	}{
		{
			name:  "Should move past tasks to today",
			today: "2024-05-01",
			initTasks: []entry.Entry{
				createEntry(1, "2024-04-29", 0, entry.StatusTaskCreated),
				createEntry(2, "2024-04-30", 0, entry.StatusTaskCreated),
			},
			wantResponse: []entry.Entry{
				createEntry(1, "2024-04-29", 0, entry.StatusTaskSnoozed),
				createEntry(2, "2024-04-30", 0, entry.StatusTaskSnoozed),
				createEntry(2, "2024-05-01", 0, entry.StatusTaskCreated),
				createEntry(1, "2024-05-01", 1, entry.StatusTaskCreated),
			},
		},
		{
			name:  "Should move past monthly tasks to this month",
			today: "2024-05",
			initTasks: []entry.Entry{
				createEntry(1, "2024-03", 0, entry.StatusTaskCreated),
				createEntry(2, "2024-04", 0, entry.StatusTaskCreated),
			},
			wantResponse: []entry.Entry{
				createEntry(1, "2024-03", 0, entry.StatusTaskSnoozed),
				createEntry(2, "2024-04", 0, entry.StatusTaskSnoozed),
				createEntry(2, "2024-05", 0, entry.StatusTaskCreated),
				createEntry(1, "2024-05", 1, entry.StatusTaskCreated),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM entries")

			var createdTasks []entry.Entry
			for _, tt := range tc.initTasks {
				db.Create(&tt)
				createdTasks = append(createdTasks, tt)
			}
			url := fmt.Sprintf("/api/tasks/pastTasks/%s/import", tc.today)
			res := api.Post(url)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}
			var resTasks []entry.Entry
			err := db.Model(entry.Entry{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to fetch Tasks from DB for comparison: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(entry.Entry{}, "ID", "TaskID")); diff != "" {
				t.Errorf("ImportPastTasks() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestImportTaskFromRecurringTask(t *testing.T) {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	es := database.NewEntryService(db)
	rts := database.NewRecurringTaskService(db)

	_, api := humatest.New(t)
	ImportPastTasksHandler(api, rts, es)

	createEntry := func(id int, title string, dateParam string, rank int, status entry.Status) entry.Entry {
		return entry.Entry{
			TaskID:          strconv.Itoa(id),
			Title:           title,
			Status:          status,
			CreatedDate:     date.DateString(dateParam),
			TaskUpdate:      "",
			Rank:            rank,
			RecurringTaskID: 1,
		}
	}

	testCases := []struct {
		name         string
		today        string
		reccTask     recurringtask.RecurringTask
		initTasks    []entry.Entry
		wantResponse []entry.Entry
	}{
		{
			name:  "recurring_task_created",
			today: "2024-05-01",
			reccTask: recurringtask.RecurringTask{
				ID:              1,
				TaskTitle:       "recurring task 1",
				TaskDescription: "",
				FreqByWeekDay:   "WED",
			},
			initTasks: []entry.Entry{},
			wantResponse: []entry.Entry{
				createEntry(1, "recurring task 1 - 2024-05-01", "2024-05-01", 0, entry.StatusTaskCreated),
			},
		},
		{
			name:  "recurring_task_different_weekday",
			today: "2024-05-01",
			reccTask: recurringtask.RecurringTask{
				ID:              1,
				TaskTitle:       "recurring task 1",
				TaskDescription: "",
				FreqByWeekDay:   "SUN",
			},
			initTasks:    []entry.Entry{},
			wantResponse: []entry.Entry{},
		},
		{
			name:  "recurring_task_already_created",
			today: "2024-05-01",
			reccTask: recurringtask.RecurringTask{
				ID:              1,
				TaskTitle:       "recurring task 1",
				TaskDescription: "",
				FreqByWeekDay:   "WED",
			},
			initTasks: []entry.Entry{
				createEntry(1, "recurring task 1 - 2024-05-01 - changed", "2024-05-01", 0, entry.StatusTaskCreated),
			},
			wantResponse: []entry.Entry{
				createEntry(1, "recurring task 1 - 2024-05-01 - changed", "2024-05-01", 0, entry.StatusTaskCreated),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM entries")
			db.Exec("DELETE FROM recurring_tasks")

			err := db.Save(&tc.reccTask).Error
			if err != nil {
				t.Fatalf("Failed to save recurring task: %v", err)
			}
			var createdTasks []entry.Entry
			for _, tt := range tc.initTasks {
				db.Create(&tt)
				createdTasks = append(createdTasks, tt)
			}

			url := fmt.Sprintf("/api/tasks/pastTasks/%s/import", tc.today)
			res := api.Post(url)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}
			var resTasks []entry.Entry
			err = db.Model(entry.Entry{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to fetch Tasks from DB for comparison: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(entry.Entry{}, "ID", "TaskID")); diff != "" {
				t.Errorf("ImportPastTasks() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
