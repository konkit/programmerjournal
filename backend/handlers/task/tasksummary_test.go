package task

import (
	"encoding/json"
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"os"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/entry"
	"strconv"
	"testing"
)

func TestGetTaskSummary(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := database.InitDB(dbTestPath)
	defer os.Remove(dbTestPath)
	es := database.NewEntryService(db)

	_, api := humatest.New(t)
	GetTaskSummaryHandler(api, es)

	testCases := []struct {
		name         string
		initTasks    []entry.Entry
		wantResponse entry.TaskSummary
		taskID       string
	}{
		{
			name:   "list single task",
			taskID: "1234",
			initTasks: []entry.Entry{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      entry.StatusTaskCreated,
					CreatedDate: "2024-05-01",
					TaskUpdate:  "TaskUpdate 1",
				},
				{
					TaskID:      "111",
					Title:       "test 2",
					Status:      entry.StatusTaskCreated,
					CreatedDate: "2024-05-01",
					TaskUpdate:  "TaskUpdate 2",
				},
				{
					TaskID:      "1234",
					Title:       "test 3",
					Status:      entry.StatusTaskCreated,
					CreatedDate: "2024-05-02",
					TaskUpdate:  "TaskUpdate 3",
				},
			},
			wantResponse: entry.TaskSummary{
				TaskEntry: entry.Entry{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      entry.StatusTaskCreated,
					CreatedDate: "2024-05-01",
					TaskUpdate:  "TaskUpdate 1",
				},
				Updates: []entry.TaskUpdate{
					{
						Date:   "2024-05-01",
						Update: "TaskUpdate 1",
						Status: entry.StatusTaskCreated,
					},
					{
						Date:   "2024-05-02",
						Update: "TaskUpdate 3",
						Status: entry.StatusTaskCreated,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM entries")

			var taskEntryIDs []string
			for _, ttt := range tc.initTasks {
				db.Create(&ttt)
				createdIDAsString := strconv.Itoa(int(ttt.ID))
				taskEntryIDs = append(taskEntryIDs, createdIDAsString)
			}

			url := fmt.Sprintf("/api/tasks/%s/summary", taskEntryIDs[0])
			res := api.Get(url)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			resTasks := entry.TaskSummary{}
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
