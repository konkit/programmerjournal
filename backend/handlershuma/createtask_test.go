package handlershuma

import (
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
	"testing"
)

func TestCreateTasks(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)

	_, api := humatest.New(t)
	CreateTask(api, dbRepo)

	testCases := []struct {
		name         string
		initTasks    []task.Task
		wantResponse []task.Task
		date         string
	}{
		{
			name:      "create a task",
			initTasks: []task.Task{},
			wantResponse: []task.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
					Rank:        0,
				},
			},
			date: "2024-05-01",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanupTaskDB(db)
			})
			for _, task := range tc.initTasks {
				db.Create(&task)
			}

			newTask := task.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      task.StatusCreated,
				CreatedDate: "2024-05-01",
				Update:      "",
			}
			//var buf bytes.Buffer
			//err := json.NewEncoder(&buf).Encode(newTask)
			//req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/create"), &buf)
			//res := httptest.NewRecorder()
			//
			//h.CreateTask(res, req)

			url := "/api/tasks/create"
			res := api.Post(url, newTask)

			if res.Code != http.StatusCreated {
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
