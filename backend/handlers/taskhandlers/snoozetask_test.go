package taskhandlers

import (
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"os"
	"programmerjournal-backend/model/entry"
	"testing"
)

func TestSnoozeTask(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := entry.InitDB(dbTestPath)
	defer os.Remove(dbTestPath)

	dbRepo, _ := entry.NewRepository(db)

	_, api := humatest.New(t)
	SnoozeTask(api, dbRepo)

	testCases := []struct {
		name         string
		initTask     entry.Entry
		wantResponse []entry.Entry
		date         string
	}{
		{
			name: "snooze task",
			initTask: entry.Entry{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      entry.StatusTaskCreated,
				CreatedDate: "2024-05-01",
				Rank:        0,
				TaskUpdate:  "",
			},
			wantResponse: []entry.Entry{
				{
					TaskID:           "1234",
					Title:            "test 1",
					Status:           entry.StatusTaskSnoozed,
					CreatedDate:      "2024-05-01",
					Rank:             0,
					TaskUpdate:       "",
					TaskSnoozedUntil: "2024-05-05",
				},
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      entry.StatusTaskCreated,
					CreatedDate: "2024-05-05",
					TaskUpdate:  "",
					Rank:        0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM tasks")

			db.Create(&tc.initTask)

			insertedID := tc.initTask.ID

			snoozeEntry := struct{ Date string }{
				Date: "2024-05-05",
			}
			url := fmt.Sprintf("/api/tasks/%d/snooze", insertedID)
			res := api.Patch(url, snoozeEntry)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status 201, got %d", res.Code)
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
