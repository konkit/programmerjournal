package handlershuma

import (
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
	"strconv"
	"testing"
)

func TestDeleteTasks(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)

	_, api := humatest.New(t)
	DeleteTask(api, dbRepo)

	testCases := []struct {
		name         string
		initTask     task.Task
		wantResponse []task.Task
		date         string
	}{
		{
			name: "delete task",
			initTask: task.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      task.StatusCreated,
				CreatedDate: "2024-05-01",
				Update:      "",
			},
			wantResponse: []task.Task{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanupTaskDB(db)
			})

			db.Create(&tc.initTask)
			taskID := tc.initTask.ID

			url := fmt.Sprintf("/api/tasks/%s/delete", strconv.Itoa(int(taskID)))
			res := api.Delete(url)

			//req := httptest.NewRequest(http.MethodGet, url, nil)
			//req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(taskID))})
			//res := httptest.NewRecorder()
			//
			//h.DeleteTask(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
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
