package handlershuma

import (
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"programmerjournal-backend/handlers"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
	"testing"
)

func TestSnoozeTask(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)

	_, api := humatest.New(t)
	SnoozeTask(api, dbRepo)

	testCases := []struct {
		name         string
		initTask     task.Task
		wantResponse []task.Task
		date         string
	}{
		{
			name: "snooze task",
			initTask: task.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      task.StatusCreated,
				CreatedDate: "2024-05-01",
				Rank:        0,
				Update:      "",
			},
			wantResponse: []task.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      task.StatusSnoozed,
					CreatedDate: "2024-05-01",
					Rank:        0,
					Update:      "",
				},
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-05",
					Update:      "",
					Rank:        0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanupTaskDB(db)
			})

			db.Create(&tc.initTask)

			insertedID := tc.initTask.ID

			snoozeEntry := handlers.SnoozeTaskEntry{
				Date: "2024-05-05",
			}
			url := fmt.Sprintf("/api/tasks/%d/snooze", insertedID)
			res := api.Patch(url, snoozeEntry)

			//var buf bytes.Buffer
			//err := json.NewEncoder(&buf).Encode(snoozeEntry)
			//req := httptest.NewRequest(http.MethodGet, url, &buf)
			//req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(insertedID))})
			//res := httptest.NewRecorder()
			//
			//h.SnoozeTask(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status 201, got %d", res.Code)
			}

			var resTasks []task.Task
			err := db.Model(task.Task{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(task.Task{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
