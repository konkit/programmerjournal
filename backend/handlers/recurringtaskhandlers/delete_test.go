package recurringtaskhandlers

import (
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"os"
	"programmerjournal-backend/model/database"
	"programmerjournal-backend/model/recurringtask"
	"strconv"
	"testing"
)

func TestDeleteTasks(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := database.InitDB(dbTestPath)
	defer os.Remove(dbTestPath)

	service := recurringtask.NewService(db)

	_, api := humatest.New(t)
	Delete(api, service)

	testCases := []struct {
		name         string
		initTask     recurringtask.RecurringTask
		wantResponse []recurringtask.RecurringTask
		date         string
	}{
		{
			name: "delete task",
			initTask: recurringtask.RecurringTask{
				TaskTitle:       "test 1",
				TaskDescription: "description 1",
				FreqByWeekDay:   "MON",
			},
			wantResponse: []recurringtask.RecurringTask{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM recurring_tasks")

			db.Create(&tc.initTask)
			url := fmt.Sprintf("/api/recurring/%s/delete", strconv.Itoa(int(tc.initTask.ID)))
			res := api.Delete(url)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
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
