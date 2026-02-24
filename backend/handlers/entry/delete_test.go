package entry

import (
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/entry"
	"strconv"
	"testing"
)

func TestDeleteTasks(t *testing.T) {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	es := database.NewEntryService(db)

	_, api := humatest.New(t)
	DeleteEntryHandler(api, es)

	testCases := []struct {
		name             string
		initTasks        []entry.Entry
		deletedTaskIndex int
		wantResponse     []entry.Entry
		date             string
	}{
		{
			name: "delete task",
			initTasks: []entry.Entry{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      entry.StatusTaskCreated,
					CreatedDate: "2024-05-01",
					TaskUpdate:  "",
				},
			},
			deletedTaskIndex: 0,
			wantResponse:     []entry.Entry{},
		},
		{
			name: "delete one of two",
			initTasks: []entry.Entry{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      entry.StatusTaskCreated,
					CreatedDate: "2024-05-01",
					TaskUpdate:  "",
					Rank:        0,
				},
				{
					TaskID:      "12",
					Title:       "test 2",
					Status:      entry.StatusTaskCreated,
					CreatedDate: "2024-05-01",
					TaskUpdate:  "",
					Rank:        1,
				},
			},
			deletedTaskIndex: 0,
			wantResponse: []entry.Entry{
				{
					TaskID:      "12",
					Title:       "test 2",
					Status:      entry.StatusTaskCreated,
					CreatedDate: "2024-05-01",
					TaskUpdate:  "",
					Rank:        0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM entries")

			for i := 0; i < len(tc.initTasks); i++ {
				db.Create(&tc.initTasks[i])
			}

			deletedTaskID := tc.initTasks[tc.deletedTaskIndex].ID
			url := fmt.Sprintf("/api/entries/%s/delete", strconv.Itoa(int(deletedTaskID)))
			res := api.Delete(url)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			var resTasks []entry.Entry
			err := db.Model(entry.Entry{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(entry.Entry{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
