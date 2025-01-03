package handlershuma

import (
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
	"strconv"
	"testing"
)

func TestImportPastTasks(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)

	_, api := humatest.New(t)
	ImportPastTasks(api, dbRepo)

	createEntry := func(id int, date string, rank int, status task.Status) task.Task {
		return task.Task{
			TaskID:      strconv.Itoa(id),
			Title:       fmt.Sprintf("test %d", id),
			Status:      status,
			CreatedDate: date,
			Update:      "",
			Rank:        rank,
		}
	}

	testCases := []struct {
		name         string
		today        string
		initTasks    []task.Task
		wantResponse []task.Task
	}{
		{
			name:  "Should move past tasks to today",
			today: "2024-05-01",
			initTasks: []task.Task{
				createEntry(1, "2024-04-21", 0, task.StatusCreated),
				createEntry(2, "2024-04-22", 0, task.StatusCreated),
			},
			wantResponse: []task.Task{
				createEntry(1, "2024-04-21", 0, task.StatusSnoozed),
				createEntry(2, "2024-04-22", 0, task.StatusSnoozed),
				createEntry(2, "2024-05-01", 0, task.StatusCreated),
				createEntry(1, "2024-05-01", 1, task.StatusCreated),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanupTaskDB(db)
			})

			var createdTasks []task.Task
			for _, tt := range tc.initTasks {
				db.Create(&tt)
				createdTasks = append(createdTasks, tt)
			}
			url := fmt.Sprintf("/api/tasks/importPastTasks/%s", tc.today)
			res := api.Post(url)

			//req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/tasks/loadDay/%s", tc.today), nil)
			//req = mux.SetURLVars(req, map[string]string{"date": tc.today})
			//res := httptest.NewRecorder()

			//h.ImportPastTasks(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			// Verify

			var resTasks []task.Task
			err := db.Model(task.Task{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to fetch Tasks from DB for comparison: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(task.Task{}, "ID", "TaskID")); diff != "" {
				t.Errorf("ImportPastTasks() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
