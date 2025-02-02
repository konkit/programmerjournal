package task

import (
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"os"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/entry"
	"testing"
)

func TestSetTaskUpdate(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := database.InitDB(dbTestPath)
	defer os.Remove(dbTestPath)
	es := database.EntryService(db)

	_, api := humatest.New(t)
	SetTaskUpdateHandler(api, es)

	testCases := []struct {
		name         string
		initTask     entry.Entry
		wantResponse entry.Entry
		date         string
	}{
		{
			name: "set task title",
			initTask: entry.Entry{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      entry.StatusTaskCreated,
				CreatedDate: "2024-05-01",
				TaskUpdate:  "",
			},
			wantResponse: entry.Entry{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      entry.StatusTaskCreated,
				CreatedDate: "2024-05-01",
				TaskUpdate:  "asdf",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM entries")

			db.Create(&tc.initTask)

			insertedID := tc.initTask.ID

			e := struct{ Update string }{
				Update: "asdf",
			}
			url := fmt.Sprintf("/api/tasks/%d/setUpdate", insertedID)
			res := api.Patch(url, e)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status 201, got %d", res.Code)
			}

			var resTasks = entry.Entry{ID: insertedID}
			db.First(&resTasks)

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(entry.Entry{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
