package taskhandlers

import (
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

func TestMigrateToDailyTask(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := database.InitDB(dbTestPath)
	defer os.Remove(dbTestPath)

	dbRepo := entry.NewService(db)

	_, api := humatest.New(t)
	MigrateTaskToDailyLog(api, dbRepo)

	testCases := []struct {
		name         string
		initTask     entry.Entry
		wantResponse []entry.Entry
		date         string
	}{
		{
			name: "migrate to Daily",
			initTask: entry.Entry{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      entry.StatusTaskCreated,
				CreatedDate: "2024-05",
				Rank:        0,
				TaskUpdate:  "",
			},
			date: "2024-05-01",
			wantResponse: []entry.Entry{
				{
					TaskID:           "1234",
					Title:            "test 1",
					Status:           entry.StatusTaskMigrated,
					CreatedDate:      "2024-05",
					Rank:             0,
					TaskUpdate:       "",
					TaskSnoozedUntil: "2024-05-01",
				},
				{
					TaskID:      "1234",
					Title:       "test 1",
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

			db.Create(&tc.initTask)

			insertedID := tc.initTask.ID

			migrateEntry := struct{ Date string }{
				Date: tc.date,
			}
			url := fmt.Sprintf("/api/tasks/%d/migrateToDaily", insertedID)
			res := api.Patch(url, migrateEntry)

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
