package handlershuma

import (
	"encoding/json"
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"gorm.io/gorm"
	"net/http"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
	"testing"
)

func TestListTasks(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)

	_, api := humatest.New(t)
	ListTasks(api, dbRepo)

	testCases := []struct {
		name         string
		initTasks    []task.Task
		wantResponse []task.Task
		date         string
	}{
		{
			name:         "empty response",
			initTasks:    []task.Task{},
			wantResponse: []task.Task{},
			date:         "2024-05-01",
		},
		{
			name: "list single task",
			initTasks: []task.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
			},
			wantResponse: []task.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
			},
			date: "2024-05-01",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanupTaskDB(db)
			})

			for _, task := range tc.initTasks {
				db.Create(&task)
			}

			url := fmt.Sprintf("/api/tasks/list/%s", "2024-05-01")
			res := api.Get(url)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			resTasks := []task.Task{}
			err := json.NewDecoder(res.Body).Decode(&resTasks)
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(task.Task{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func cleanupTaskDB(db *gorm.DB) *gorm.DB {
	return db.Exec("DELETE FROM tasks")
}
