package entry

import (
	"encoding/json"
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"os"
	"programmerjournal-backend/model/database"
	"programmerjournal-backend/model/entry"
	"testing"
)

func TestListTasks(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := database.InitDB(dbTestPath)
	defer os.Remove(dbTestPath)

	_, api := humatest.New(t)
	ListEntriesHandler(api, db)

	testCases := []struct {
		name         string
		initTasks    []entry.Entry
		wantResponse []entry.Entry
		date         string
	}{
		{
			name:         "empty response",
			initTasks:    []entry.Entry{},
			wantResponse: []entry.Entry{},
			date:         "2024-05-01",
		},
		{
			name: "list single task",
			initTasks: []entry.Entry{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      entry.StatusTaskCreated,
					CreatedDate: "2024-05-01",
					TaskUpdate:  "",
				},
			},
			wantResponse: []entry.Entry{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      entry.StatusTaskCreated,
					CreatedDate: "2024-05-01",
					TaskUpdate:  "",
				},
			},
			date: "2024-05-01",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM entries")

			for _, task := range tc.initTasks {
				db.Create(&task)
			}

			url := fmt.Sprintf("/api/entries/list/%s", "2024-05-01")
			res := api.Get(url)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			resTasks := []entry.Entry{}
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
