package handlers

import (
	"encoding/json"
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

func TestChangeRank(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)

	_, api := humatest.New(t)
	ChangeTaskRank(api, dbRepo)
	ListTasks(api, dbRepo)

	createEntry := func(titleID int, rank int) task.Task {
		return task.Task{
			TaskID:      strconv.Itoa(titleID),
			Title:       fmt.Sprintf("test %d", titleID),
			Status:      task.StatusCreated,
			CreatedDate: "2024-05-01",
			Update:      "",
			Rank:        rank,
		}
	}

	testCases := []struct {
		name         string
		initTasks    []task.Task
		oldRank      int
		newRank      int
		wantResponse []task.Task
	}{
		{
			name: "change rank in the middle",
			initTasks: []task.Task{
				createEntry(0, 0),
				createEntry(1, 1),
				createEntry(2, 2),
				createEntry(3, 3),
			},
			oldRank: 1,
			newRank: 3,
			wantResponse: []task.Task{
				createEntry(0, 0),
				createEntry(2, 1),
				createEntry(3, 2),
				createEntry(1, 3),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanupTaskDB(db)
			})

			var createdTasks []task.Task
			for _, tt := range tc.initTasks {
				db.Create(&tt)
				createdTasks = append(createdTasks, tt)
			}

			changedID := createdTasks[tc.oldRank].ID
			entry := struct{ NewRank int }{
				NewRank: tc.newRank,
			}
			url := fmt.Sprintf("/api/tasks/%d/changeRank", changedID)
			res := api.Patch(url, entry)

			//var buf bytes.Buffer
			//err := json.NewEncoder(&buf).Encode(entry)
			//if err != nil {
			//	t.Fatalf("error decoding entry: %v", err)
			//}
			//req := httptest.NewRequest(http.MethodPost, url, &buf)
			//req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(changedID))})
			//res := httptest.NewRecorder()
			//
			//h.ChangeRank(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			//req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/list/%s", "2024-05-01"), nil)
			//req = mux.SetURLVars(req, map[string]string{"date": "2024-05-01"})
			//res = httptest.NewRecorder()
			//
			//h.ListTasks(res, req)

			listUrl := fmt.Sprintf("/api/tasks/list/%s", "2024-05-01")
			res = api.Get(listUrl)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			resTasks := []task.Task{}
			err := json.NewDecoder(res.Body).Decode(&resTasks)
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(task.Task{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
