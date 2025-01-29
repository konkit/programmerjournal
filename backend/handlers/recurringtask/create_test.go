package recurringtask

import (
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"os"
	"programmerjournal-backend/model/database"
	"programmerjournal-backend/model/recurringtask"
	"testing"
)

func TestCreateRecurringTask(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := database.InitDB(dbTestPath)
	defer os.Remove(dbTestPath)

	_, api := humatest.New(t)
	CreateHandler(api, db)

	testCases := []struct {
		name         string
		initTasks    []recurringtask.RecurringTask
		input        CreateRecurringTaskInputBody
		wantResponse []recurringtask.RecurringTask
		date         string
	}{
		{
			name:      "create_recurring_task",
			initTasks: []recurringtask.RecurringTask{},
			input: CreateRecurringTaskInputBody{
				TaskTitle:       "recurring task 1",
				TaskDescription: "description 1",
				FreqByWeekDay:   "MON",
			},
			wantResponse: []recurringtask.RecurringTask{
				{
					TaskTitle:       "recurring task 1",
					TaskDescription: "description 1",
					FreqByWeekDay:   "MON",
				},
			},
			date: "2024-05-01",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM recurring_tasks")

			for _, task := range tc.initTasks {
				db.Create(&task)
			}

			url := "/api/recurring/create"
			res := api.Post(url, tc.input)

			if res.Code != http.StatusCreated {
				t.Fatalf("Expected status 201, got %d", res.Code)
			}

			var resTasks []recurringtask.RecurringTask
			err := db.Model(recurringtask.RecurringTask{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(recurringtask.RecurringTask{}, "ID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
