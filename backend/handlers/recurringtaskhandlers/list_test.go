package recurringtaskhandlers

import (
	"encoding/json"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"os"
	"programmerjournal-backend/model/database"
	"programmerjournal-backend/model/recurringtask"
	"testing"
)

func TestListRecurringTasks(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := database.InitDB(dbTestPath)
	defer os.Remove(dbTestPath)

	service := recurringtask.NewService(db)

	_, api := humatest.New(t)
	List(api, service)

	testCases := []struct {
		name         string
		initTasks    []recurringtask.RecurringTask
		wantResponse []recurringtask.RecurringTask
		date         string
	}{
		{
			name:         "empty response",
			initTasks:    []recurringtask.RecurringTask{},
			wantResponse: []recurringtask.RecurringTask{},
		},
		{
			name: "list single recurring task",
			initTasks: []recurringtask.RecurringTask{
				{
					TaskTitle:       "test 1",
					TaskDescription: "description 1",
					FreqByWeekDay:   "MON",
				},
			},
			wantResponse: []recurringtask.RecurringTask{
				{
					TaskTitle:       "test 1",
					TaskDescription: "description 1",
					FreqByWeekDay:   "MON",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM recurring_tasks")

			for _, task := range tc.initTasks {
				db.Create(&task)
			}

			url := "/api/recurring/list"
			res := api.Get(url)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			resTasks := []recurringtask.RecurringTask{}
			err := json.NewDecoder(res.Body).Decode(&resTasks)
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(recurringtask.RecurringTask{}, "ID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
