package handlers

import (
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
	"testing"
)

func TestSetTaskTitle(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)

	_, api := humatest.New(t)
	SetTaskTitle(api, dbRepo)

	testCases := []struct {
		name         string
		initTask     task.Task
		wantResponse task.Task
		date         string
	}{
		{
			name: "set task title",
			initTask: task.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      task.StatusCreated,
				CreatedDate: "2024-05-01",
				Update:      "",
			},
			wantResponse: task.Task{
				TaskID:      "1234",
				Title:       "test 2",
				Status:      task.StatusCreated,
				CreatedDate: "2024-05-01",
				Update:      "",
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

			entry := struct{ Title string }{
				Title: "test 2",
			}
			url := fmt.Sprintf("/api/tasks/%d/setTitle", insertedID)
			res := api.Patch(url, entry)

			//var buf bytes.Buffer
			//err := json.NewEncoder(&buf).Encode(entry)
			//if err != nil {
			//	t.Fatalf("error decoding entry: %v", err)
			//}
			//req := httptest.NewRequest(http.MethodPost, url, &buf)
			//req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(insertedID))})
			//res := httptest.NewRecorder()

			//h.SetTaskTitle(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status 201, got %d", res.Code)
			}

			var resTasks = task.Task{ID: insertedID}
			db.First(&resTasks)

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(task.Task{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
