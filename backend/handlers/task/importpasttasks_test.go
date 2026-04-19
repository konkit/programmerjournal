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

func TestImportWeeklyMigratedTask(t *testing.T) {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	es := database.NewEntryService(db)
	rts := database.NewRecurringTaskService(db)

	_, api := humatest.New(t)
	ImportPastTasksHandler(api, rts, es)

	// 1. Create a weekly task for 2024-W18 (which contains 2024-05-01)
	weeklyTask := entry.Entry{
		TaskID:           "weekly-1",
		Title:            "Weekly Task",
		Status:           entry.StatusTaskSnoozed,
		CreatedDate:      "2024-W18",
		TaskSnoozedUntil: "2024-04-30",
	}
	db.Create(&weeklyTask)

	// 2. Create the migrated task on 2024-04-30
	migratedTask := entry.Entry{
		TaskID:      "weekly-1",
		Title:       "Weekly Task",
		Status:      entry.StatusTaskCreated,
		CreatedDate: "2024-04-30",
	}
	db.Create(&migratedTask)

	// 3. Import tasks to 2024-05-01 (May 1st 2024 is Wednesday, Week 18)
	url := "/api/tasks/pastTasks/2024-05-01/import"
	res := api.Post(url)

	if res.Code != http.StatusOK {
		t.Fatalf("Expected status OK, got %d", res.Code)
	}

	// 4. Verify results
	var resTasks []entry.Entry
	err = db.Model(entry.Entry{}).Find(&resTasks).Error
	if err != nil {
		t.Fatalf("Failed to fetch Tasks from DB for comparison: %v", err)
	}

	// We expect:
	// - Weekly task to be StatusTaskCreated on 2024-W18
	// - Migrated task on 2024-04-30 to be StatusTaskSnoozed
	// - NO NEW task on 2024-05-01

	for _, rt := range resTasks {
		if rt.CreatedDate == "2024-W18" {
			if rt.Status != entry.StatusTaskCreated {
				t.Errorf("Expected weekly task to be StatusTaskCreated, got %s", rt.Status)
			}
		} else if rt.CreatedDate == "2024-04-30" {
			if rt.Status != entry.StatusTaskSnoozed {
				t.Errorf("Expected past task to be StatusTaskSnoozed, got %s", rt.Status)
			}
		} else if rt.CreatedDate == "2024-05-01" {
			t.Errorf("Expected NO task on 2024-05-01, but found one: %s", rt.Title)
		}
	}
}

func TestImportWeeklyMigratedSnoozeChainTask(t *testing.T) {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	es := database.NewEntryService(db)
	rts := database.NewRecurringTaskService(db)

	_, api := humatest.New(t)
	ImportPastTasksHandler(api, rts, es)

	// Chain: Weekly -> 2024-04-29 (Snoozed) -> 2024-04-30 (Created)
	// Target: 2024-05-01

	weeklyTask := entry.Entry{
		TaskID:           "weekly-chain",
		Title:            "Weekly Chain Task",
		Status:           entry.StatusTaskSnoozed,
		CreatedDate:      "2024-W18",
		TaskSnoozedUntil: "2024-04-29",
	}
	db.Create(&weeklyTask)

	day1Task := entry.Entry{
		TaskID:           "weekly-chain",
		Title:            "Weekly Chain Task",
		Status:           entry.StatusTaskSnoozed,
		CreatedDate:      "2024-04-29",
		TaskSnoozedUntil: "2024-04-30",
	}
	db.Create(&day1Task)

	day2Task := entry.Entry{
		TaskID:      "weekly-chain",
		Title:       "Weekly Chain Task",
		Status:      entry.StatusTaskCreated,
		CreatedDate: "2024-04-30",
	}
	db.Create(&day2Task)

	url := "/api/tasks/pastTasks/2024-05-01/import"
	res := api.Post(url)

	if res.Code != http.StatusOK {
		t.Fatalf("Expected status OK, got %d", res.Code)
	}

	var resTasks []entry.Entry
	err = db.Model(entry.Entry{}).Find(&resTasks).Error
	if err != nil {
		t.Fatalf("Failed to fetch Tasks from DB for comparison: %v", err)
	}

	for _, rt := range resTasks {
		if rt.CreatedDate == "2024-W18" {
			if rt.Status != entry.StatusTaskCreated {
				t.Errorf("Expected weekly task to be StatusTaskCreated, got %s", rt.Status)
			}
		} else if rt.TaskID == "weekly-chain" && date.GetDateType(string(rt.CreatedDate)) == date.DateTypeDay {
			if rt.Status != entry.StatusTaskSnoozed {
				t.Errorf("Expected daily task on %s to be StatusTaskSnoozed, got %s", rt.CreatedDate, rt.Status)
			}
		} else if rt.CreatedDate == "2024-05-01" {
			t.Errorf("Expected NO task on 2024-05-01, but found one: %s", rt.Title)
		}
	}
}
