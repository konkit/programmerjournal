package taskhandlers

import (
	"encoding/json"
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"os"
	"programmerjournal-backend/model/entry"
	"testing"
)

func TestWeeklySummary(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := entry.InitDB(dbTestPath)
	defer os.Remove(dbTestPath)

	dbRepo, _ := entry.NewRepository(db)

	_, api := humatest.New(t)
	WeeklyTaskSummary(api, dbRepo)

	testCases := []struct {
		name         string
		initTasks    []entry.Entry
		wantResponse []entry.TaskSummary
		date         string
	}{
		{
			name:         "empty_response",
			initTasks:    []entry.Entry{},
			wantResponse: nil,
			date:         "2024-05-01",
		},
		{
			name: "single_task_two_updates",
			initTasks: []entry.Entry{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      entry.StatusTaskCreated,
					CreatedDate: "2024-05-01",
					TaskUpdate:  "Wednesday update",
				},
				{
					TaskID:      "1234",
					Title:       "test 1 - updated",
					Status:      entry.StatusTaskCreated,
					CreatedDate: "2024-05-02",
					TaskUpdate:  "Thursday update",
				},
			},
			wantResponse: []entry.TaskSummary{
				{
					TaskEntry: entry.Entry{
						TaskID:      "1234",
						Title:       "test 1 - updated",
						Status:      entry.StatusTaskCreated,
						CreatedDate: "2024-05-02",
						TaskUpdate:  "Thursday update",
					},
					Updates: []entry.TaskUpdate{
						{
							Date:   "2024-05-01",
							Update: "Wednesday update",
							Status: entry.StatusTaskCreated,
						},
						{
							Date:   "2024-05-02",
							Update: "Thursday update",
							Status: entry.StatusTaskCreated,
						},
					},
				},
			},
			date: "2024-04-29",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM tasks")

			for _, task := range tc.initTasks {
				db.Create(&task)
			}

			url := fmt.Sprintf("/api/tasks/weeklySummary/%s", "2024-04-29")
			res := api.Get(url)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			resTasks := []entry.TaskSummary{}
			err := json.NewDecoder(res.Body).Decode(&resTasks)
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(entry.Entry{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
