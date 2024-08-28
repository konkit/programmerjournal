package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"net/http/httptest"
	"programmerjournal-backend/handlers"
	"strconv"
	"testing"
)

const dbTestPath = "./test.db"

func TestListTasks(t *testing.T) {
	db := handlers.InitDB(dbTestPath)
	h := handlers.NewHandlers(db)

	testCases := []struct {
		name         string
		initTasks    []handlers.Task
		wantResponse []handlers.Task
		date         string
	}{
		{
			name:         "empty response",
			initTasks:    []handlers.Task{},
			wantResponse: []handlers.Task{},
			date:         "2024-05-01",
		},
		{
			name: "list single task",
			initTasks: []handlers.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      handlers.TaskCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
			},
			wantResponse: []handlers.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      handlers.TaskCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
			},
			date: "2024-05-01",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, task := range tc.initTasks {
				db.Create(&task)
			}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/list/%s", "2024-05-01"), nil)
			req = mux.SetURLVars(req, map[string]string{"date": tc.date})
			res := httptest.NewRecorder()

			h.ListTasks(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			resTasks := []handlers.Task{}
			err := json.NewDecoder(res.Body).Decode(&resTasks)
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(handlers.Task{}, "ID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}

			db.Exec("DELETE FROM tasks")
		})
	}
}

func TestCreateTasks(t *testing.T) {
	db := handlers.InitDB(dbTestPath)
	h := handlers.NewHandlers(db)

	testCases := []struct {
		name         string
		initTasks    []handlers.Task
		wantResponse []handlers.Task
		date         string
	}{
		{
			name:      "create a task",
			initTasks: []handlers.Task{},
			wantResponse: []handlers.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      handlers.TaskCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
			},
			date: "2024-05-01",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, task := range tc.initTasks {
				db.Create(&task)
			}

			newTask := handlers.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      handlers.TaskCreated,
				CreatedDate: "2024-05-01",
				Update:      "",
			}
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(newTask)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/create"), &buf)
			res := httptest.NewRecorder()

			h.CreateTask(res, req)

			if res.Code != http.StatusCreated {
				t.Fatalf("Expected status 201, got %d", res.Code)
			}

			var resTasks []handlers.Task
			err = db.Model(handlers.Task{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(handlers.Task{}, "ID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}

			db.Exec("DELETE FROM tasks")
		})
	}
}

func TestUpdateTasks(t *testing.T) {
	db := handlers.InitDB(dbTestPath)
	h := handlers.NewHandlers(db)

	testCases := []struct {
		name         string
		initTask     handlers.Task
		wantResponse []handlers.Task
		date         string
	}{
		{
			name: "update task",
			initTask: handlers.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      handlers.TaskCreated,
				CreatedDate: "2024-05-01",
				Update:      "",
			},
			wantResponse: []handlers.Task{
				{
					TaskID:      "1234",
					Title:       "test 2",
					Status:      handlers.TaskCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
			},
			date: "2024-05-01",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Create(&tc.initTask)

			insertedID := tc.initTask.ID

			updatedTask := handlers.Task{
				ID:          insertedID,
				TaskID:      "1234",
				Title:       "test 2",
				Status:      handlers.TaskCreated,
				CreatedDate: "2024-05-01",
				Update:      "",
			}
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(updatedTask)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/%d/update", insertedID), &buf)
			req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(insertedID))})
			res := httptest.NewRecorder()

			h.UpdateTask(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status 201, got %d", res.Code)
			}

			var resTasks []handlers.Task
			err = db.Model(handlers.Task{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(handlers.Task{}, "ID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}

			db.Exec("DELETE FROM tasks")
		})
	}
}

func TestDeleteTasks(t *testing.T) {
	db := handlers.InitDB(dbTestPath)
	h := handlers.NewHandlers(db)

	testCases := []struct {
		name         string
		initTask     handlers.Task
		wantResponse []handlers.Task
		date         string
	}{
		{
			name: "delete task",
			initTask: handlers.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      handlers.TaskCreated,
				CreatedDate: "2024-05-01",
				Update:      "",
			},
			wantResponse: []handlers.Task{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Create(&tc.initTask)
			taskID := tc.initTask.ID

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/delete/%s", strconv.Itoa(int(taskID))), nil)
			req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(taskID))})
			res := httptest.NewRecorder()

			h.DeleteTask(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			var resTasks []handlers.Task
			err := db.Model(handlers.Task{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(handlers.Task{}, "ID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}

			db.Exec("DELETE FROM tasks")
		})
	}
}

func TestSnoozeTask(t *testing.T) {
	db := handlers.InitDB(dbTestPath)
	h := handlers.NewHandlers(db)

	testCases := []struct {
		name         string
		initTask     handlers.Task
		wantResponse []handlers.Task
		date         string
	}{
		{
			name: "snooze task",
			initTask: handlers.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      handlers.TaskCreated,
				CreatedDate: "2024-05-01",
				Update:      "",
			},
			wantResponse: []handlers.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      handlers.TaskSnoozed,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      handlers.TaskCreated,
					CreatedDate: "2024-05-05",
					Update:      "",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Create(&tc.initTask)

			insertedID := tc.initTask.ID

			snoozeEntry := handlers.SnoozeTaskEntry{
				Date: "2024-05-05",
			}
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(snoozeEntry)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/%d/snooze", insertedID), &buf)
			req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(insertedID))})
			res := httptest.NewRecorder()

			h.SnoozeTask(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status 201, got %d", res.Code)
			}

			var resTasks []handlers.Task
			err = db.Model(handlers.Task{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(handlers.Task{}, "ID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}

			db.Exec("DELETE FROM tasks")
		})
	}
}

func TestSetTaskDone(t *testing.T) {
	db := handlers.InitDB(dbTestPath)
	h := handlers.NewHandlers(db)

	testCases := []struct {
		name         string
		initTask     handlers.Task
		wantResponse handlers.Task
		date         string
	}{
		{
			name: "set task done",
			initTask: handlers.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      handlers.TaskCreated,
				CreatedDate: "2024-05-01",
				Update:      "",
			},
			wantResponse: handlers.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      handlers.TaskDone,
				CreatedDate: "2024-05-01",
				Update:      "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Create(&tc.initTask)

			insertedID := tc.initTask.ID

			entry := handlers.SetTaskDoneEntry{
				Done: true,
			}
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(entry)
			if err != nil {
				t.Fatalf("error decoding entry: %v", err)
			}
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/tasks/%d/setDone", insertedID), &buf)
			req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(insertedID))})
			res := httptest.NewRecorder()

			h.SetTaskDone(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status 201, got %d", res.Code)
			}

			var resTasks = handlers.Task{ID: insertedID}
			db.First(&resTasks)

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(handlers.Task{}, "ID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}

			db.Exec("DELETE FROM tasks")
		})
	}
}

func TestSetTaskTitle(t *testing.T) {
	db := handlers.InitDB(dbTestPath)
	h := handlers.NewHandlers(db)

	testCases := []struct {
		name         string
		initTask     handlers.Task
		wantResponse handlers.Task
		date         string
	}{
		{
			name: "set task title",
			initTask: handlers.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      handlers.TaskCreated,
				CreatedDate: "2024-05-01",
				Update:      "",
			},
			wantResponse: handlers.Task{
				TaskID:      "1234",
				Title:       "test 2",
				Status:      handlers.TaskCreated,
				CreatedDate: "2024-05-01",
				Update:      "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Create(&tc.initTask)

			insertedID := tc.initTask.ID

			entry := handlers.SetTaskTitleEntry{
				Title: "test 2",
			}
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(entry)
			if err != nil {
				t.Fatalf("error decoding entry: %v", err)
			}
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/tasks/%d/setTitle", insertedID), &buf)
			req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(insertedID))})
			res := httptest.NewRecorder()

			h.SetTaskTitle(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status 201, got %d", res.Code)
			}

			var resTasks = handlers.Task{ID: insertedID}
			db.First(&resTasks)

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(handlers.Task{}, "ID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}

			db.Exec("DELETE FROM tasks")
		})
	}
}

func TestSetTaskUpdate(t *testing.T) {
	db := handlers.InitDB(dbTestPath)
	h := handlers.NewHandlers(db)

	testCases := []struct {
		name         string
		initTask     handlers.Task
		wantResponse handlers.Task
		date         string
	}{
		{
			name: "set task title",
			initTask: handlers.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      handlers.TaskCreated,
				CreatedDate: "2024-05-01",
				Update:      "",
			},
			wantResponse: handlers.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      handlers.TaskCreated,
				CreatedDate: "2024-05-01",
				Update:      "asdf",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Create(&tc.initTask)

			insertedID := tc.initTask.ID

			entry := handlers.SetTaskUpdateEntry{
				Update: "asdf",
			}
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(entry)
			if err != nil {
				t.Fatalf("error decoding entry: %v", err)
			}
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/tasks/%d/setUpdate", insertedID), &buf)
			req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(insertedID))})
			res := httptest.NewRecorder()

			h.SetTaskUpdate(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status 201, got %d", res.Code)
			}

			var resTasks = handlers.Task{ID: insertedID}
			db.First(&resTasks)

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(handlers.Task{}, "ID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}

			db.Exec("DELETE FROM tasks")
		})
	}
}

func TestMoveTasksToNextDay(t *testing.T) {
	db := handlers.InitDB(dbTestPath)
	h := handlers.NewHandlers(db)

	testCases := []struct {
		name         string
		currentDate  string
		nextDate     string
		initTasks    []handlers.Task
		wantResponse []handlers.Task
		date         string
	}{
		{
			name:         "empty response",
			currentDate:  "2024-05-01",
			nextDate:     "2024-05-02",
			initTasks:    []handlers.Task{},
			wantResponse: []handlers.Task{},
			date:         "2024-05-01",
		},
		{
			name:        "list single task",
			currentDate: "2024-05-01",
			nextDate:    "2024-05-02",
			initTasks: []handlers.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      handlers.TaskCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
			},
			wantResponse: []handlers.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      handlers.TaskSnoozed,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      handlers.TaskCreated,
					CreatedDate: "2024-05-02",
					Update:      "",
				},
			},
			date: "2024-05-01",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, task := range tc.initTasks {
				db.Create(&task)
			}

			entry := handlers.MoveTasksToNextDayEntry{
				Current: "2024-05-01",
				Next:    "2024-05-02",
			}
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(entry)
			if err != nil {
				t.Fatalf("error decoding entry: %v", err)
			}
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/moveToNextDay"), &buf)
			req = mux.SetURLVars(req, map[string]string{"date": tc.date})
			res := httptest.NewRecorder()

			h.MoveTasksToNextDay(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			var resTasks []handlers.Task
			err = db.Model(handlers.Task{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to fetch Tasks from DB for comparison: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(handlers.Task{}, "ID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}

			db.Exec("DELETE FROM tasks")
		})
	}
}

func TestGetTask(t *testing.T) {
	db := handlers.InitDB(dbTestPath)
	h := handlers.NewHandlers(db)

	testCases := []struct {
		name         string
		initTasks    []handlers.Task
		wantResponse []handlers.Task
		taskID       string
	}{
		{
			name:         "empty response",
			taskID:       "1234",
			initTasks:    []handlers.Task{},
			wantResponse: []handlers.Task{},
		},
		{
			name:   "list single task",
			taskID: "1234",
			initTasks: []handlers.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      handlers.TaskCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
				{
					TaskID:      "111",
					Title:       "test 2",
					Status:      handlers.TaskCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
				{
					TaskID:      "1234",
					Title:       "test 3",
					Status:      handlers.TaskCreated,
					CreatedDate: "2024-05-02",
					Update:      "",
				},
			},
			wantResponse: []handlers.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      handlers.TaskCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
				{
					TaskID:      "1234",
					Title:       "test 3",
					Status:      handlers.TaskCreated,
					CreatedDate: "2024-05-02",
					Update:      "",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, task := range tc.initTasks {
				db.Create(&task)
			}

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/summary/%s", "1234"), nil)
			req = mux.SetURLVars(req, map[string]string{"taskID": tc.taskID})
			res := httptest.NewRecorder()

			h.GetTask(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			resTasks := []handlers.Task{}
			err := json.NewDecoder(res.Body).Decode(&resTasks)
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(handlers.Task{}, "ID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}

			db.Exec("DELETE FROM tasks")
		})
	}
}
