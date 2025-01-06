package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
	"testing"
)

func TestWeeklySummary(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)

	_, api := humatest.New(t)
	WeeklySummary(api, dbRepo)

	testCases := []struct {
		name         string
		initTasks    []task.Task
		wantResponse []task.TaskWeeklySummary
		date         string
	}{
		{
			name:         "empty_response",
			initTasks:    []task.Task{},
			wantResponse: nil,
			date:         "2024-05-01",
		},
		{
			name: "single_task_two_updates",
			initTasks: []task.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-01",
					Update:      "Wednesday update",
				},
				{
					TaskID:      "1234",
					Title:       "test 1 - updated",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-02",
					Update:      "Thursday update",
				},
			},
			wantResponse: []task.TaskWeeklySummary{
				{
					Task: task.Task{
						TaskID:      "1234",
						Title:       "test 1 - updated",
						Status:      task.StatusCreated,
						CreatedDate: "2024-05-02",
						Update:      "Thursday update",
					},
					Updates: []task.TaskUpdate{
						{
							Date:   "2024-05-01",
							Update: "Wednesday update",
							Status: task.StatusCreated,
						},
						{
							Date:   "2024-05-02",
							Update: "Thursday update",
							Status: task.StatusCreated,
						},
					},
				},
			},
			date: "2024-04-29",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				err := db.Exec("DELETE FROM tasks").Error
				if err != nil {
					t.Error(err)
				}
			})

			for _, task := range tc.initTasks {
				db.Create(&task)
			}

			url := fmt.Sprintf("/api/tasks/weeklySummary/%s", "2024-04-29")
			res := api.Get(url)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			resTasks := []task.TaskWeeklySummary{}
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
