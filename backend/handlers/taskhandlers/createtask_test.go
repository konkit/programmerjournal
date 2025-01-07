package taskhandlers

import (
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"os"
	"programmerjournal-backend/model/entry"
	"testing"
)

func TestCreateTasks(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := entry.InitDB(dbTestPath)
	defer os.Remove(dbTestPath)

	dbRepo, _ := entry.NewRepository(db)

	_, api := humatest.New(t)
	CreateTask(api, dbRepo)

	testCases := []struct {
		name         string
		initTasks    []entry.Entry
		wantResponse []entry.Entry
		date         string
	}{
		{
			name:      "create a task",
			initTasks: []entry.Entry{},
			wantResponse: []entry.Entry{
				{
					//TaskID:      "1234",
					Title:       "test 1",
					Status:      entry.StatusTaskCreated,
					CreatedDate: "2024-05-01",
					TaskUpdate:  "",
					Rank:        0,
				},
			},
			date: "2024-05-01",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM tasks")

			for _, task := range tc.initTasks {
				db.Create(&task)
			}

			newTask := struct {
				Title       string
				CreatedDate string
			}{
				Title:       "test 1",
				CreatedDate: "2024-05-01",
			}

			url := "/api/tasks/create"
			res := api.Post(url, newTask)

			if res.Code != http.StatusCreated {
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
