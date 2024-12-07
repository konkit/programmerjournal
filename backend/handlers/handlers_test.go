package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"programmerjournal-backend/handlers"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
	"strconv"
	"testing"
)

const dbTestPath = "./test.db"

func TestListTasks(t *testing.T) {
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)
	h := handlers.NewHandlers(dbRepo)

	testCases := []struct {
		name         string
		initTasks    []task.Task
		wantResponse []task.Task
		date         string
	}{
		{
			name:         "empty response",
			initTasks:    []task.Task{},
			wantResponse: []task.Task{},
			date:         "2024-05-01",
		},
		{
			name: "list single task",
			initTasks: []task.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
			},
			wantResponse: []task.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/list/%s", "2024-05-01"), nil)
			req = mux.SetURLVars(req, map[string]string{"date": tc.date})
			res := httptest.NewRecorder()

			h.ListTasks(res, req)

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

func TestCreateTasks(t *testing.T) {
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)
	h := handlers.NewHandlers(dbRepo)

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
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(newTask)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/create"), &buf)
			res := httptest.NewRecorder()

			h.CreateTask(res, req)

			if res.Code != http.StatusCreated {
				t.Fatalf("Expected status 201, got %d", res.Code)
			}

			var resTasks []task.Task
			err = db.Model(task.Task{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(task.Task{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestUpdateTasks(t *testing.T) {
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)
	h := handlers.NewHandlers(dbRepo)

	testCases := []struct {
		name         string
		initTask     task.Task
		wantResponse []task.Task
		date         string
	}{
		{
			name: "update task",
			initTask: task.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      task.StatusCreated,
				CreatedDate: "2024-05-01",
				Update:      "",
			},
			wantResponse: []task.Task{
				{
					TaskID:      "1234",
					Title:       "test 2",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
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

			db.Create(&tc.initTask)

			insertedID := tc.initTask.ID

			updatedTask := task.Task{
				ID:          insertedID,
				TaskID:      "1234",
				Title:       "test 2",
				Status:      task.StatusCreated,
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

			var resTasks []task.Task
			err = db.Model(task.Task{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(task.Task{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDeleteTasks(t *testing.T) {
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)
	h := handlers.NewHandlers(dbRepo)

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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/delete/%s", strconv.Itoa(int(taskID))), nil)
			req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(taskID))})
			res := httptest.NewRecorder()

			h.DeleteTask(res, req)

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

func TestSnoozeTask(t *testing.T) {
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)
	h := handlers.NewHandlers(dbRepo)

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
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(snoozeEntry)
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/%d/snooze", insertedID), &buf)
			req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(insertedID))})
			res := httptest.NewRecorder()

			h.SnoozeTask(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status 201, got %d", res.Code)
			}

			var resTasks []task.Task
			err = db.Model(task.Task{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(task.Task{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSetTaskDone(t *testing.T) {
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)
	h := handlers.NewHandlers(dbRepo)

	testCases := []struct {
		name         string
		initTask     task.Task
		wantResponse task.Task
		date         string
	}{
		{
			name: "set task done",
			initTask: task.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      task.StatusCreated,
				CreatedDate: "2024-05-01",
				Update:      "",
			},
			wantResponse: task.Task{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      task.StatusDone,
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

			var resTasks = task.Task{ID: insertedID}
			db.First(&resTasks)

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(task.Task{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSetTaskTitle(t *testing.T) {
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)
	h := handlers.NewHandlers(dbRepo)

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

			var resTasks = task.Task{ID: insertedID}
			db.First(&resTasks)

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(task.Task{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSetTaskUpdate(t *testing.T) {
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)
	h := handlers.NewHandlers(dbRepo)

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
				Title:       "test 1",
				Status:      task.StatusCreated,
				CreatedDate: "2024-05-01",
				Update:      "asdf",
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

			var resTasks = task.Task{ID: insertedID}
			db.First(&resTasks)

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(task.Task{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestMoveTasksToNextDay(t *testing.T) {
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)
	h := handlers.NewHandlers(dbRepo)

	testCases := []struct {
		name         string
		currentDate  string
		nextDate     string
		initTasks    []task.Task
		wantResponse []task.Task
		date         string
	}{
		{
			name:         "empty response",
			currentDate:  "2024-05-01",
			nextDate:     "2024-05-02",
			initTasks:    []task.Task{},
			wantResponse: []task.Task{},
			date:         "2024-05-01",
		},
		{
			name:        "list single task",
			currentDate: "2024-05-01",
			nextDate:    "2024-05-02",
			initTasks: []task.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
			},
			wantResponse: []task.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      task.StatusSnoozed,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-02",
					Update:      "",
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

			var resTasks []task.Task
			err = db.Model(task.Task{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to fetch Tasks from DB for comparison: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(task.Task{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetTask(t *testing.T) {
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)
	h := handlers.NewHandlers(dbRepo)

	testCases := []struct {
		name         string
		initTasks    []task.Task
		wantResponse []task.Task
		taskID       string
	}{
		{
			name:         "empty response",
			taskID:       "1234",
			initTasks:    []task.Task{},
			wantResponse: []task.Task{},
		},
		{
			name:   "list single task",
			taskID: "1234",
			initTasks: []task.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
				{
					TaskID:      "111",
					Title:       "test 2",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
				{
					TaskID:      "1234",
					Title:       "test 3",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-02",
					Update:      "",
				},
			},
			wantResponse: []task.Task{
				{
					TaskID:      "1234",
					Title:       "test 1",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-01",
					Update:      "",
				},
				{
					TaskID:      "1234",
					Title:       "test 3",
					Status:      task.StatusCreated,
					CreatedDate: "2024-05-02",
					Update:      "",
				},
			},
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/summary/%s", "1234"), nil)
			req = mux.SetURLVars(req, map[string]string{"taskID": tc.taskID})
			res := httptest.NewRecorder()

			h.GetTask(res, req)

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

func TestChangeRank(t *testing.T) {
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)
	h := handlers.NewHandlers(dbRepo)

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
			entry := handlers.ChangeRankEntry{
				NewRank: tc.newRank,
			}
			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(entry)
			if err != nil {
				t.Fatalf("error decoding entry: %v", err)
			}
			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/tasks/%d/changeRank", changedID), &buf)
			req = mux.SetURLVars(req, map[string]string{"id": strconv.Itoa(int(changedID))})
			res := httptest.NewRecorder()

			h.ChangeRank(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/tasks/list/%s", "2024-05-01"), nil)
			req = mux.SetURLVars(req, map[string]string{"date": "2024-05-01"})
			res = httptest.NewRecorder()

			h.ListTasks(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			resTasks := []task.Task{}
			err = json.NewDecoder(res.Body).Decode(&resTasks)
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(task.Task{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestImportPastTasks(t *testing.T) {
	db, _ := taskrepository.InitDB(dbTestPath)
	dbRepo, _ := taskrepository.NewRepository(db)
	h := handlers.NewHandlers(dbRepo)

	createEntry := func(id int, date string, rank int, status task.Status) task.Task {
		return task.Task{
			TaskID:      strconv.Itoa(id),
			Title:       fmt.Sprintf("test %d", id),
			Status:      status,
			CreatedDate: date,
			Update:      "",
			Rank:        rank,
		}
	}

	testCases := []struct {
		name         string
		today        string
		initTasks    []task.Task
		wantResponse []task.Task
	}{
		{
			name:  "Should move past tasks to today",
			today: "2024-05-01",
			initTasks: []task.Task{
				createEntry(1, "2024-04-21", 0, task.StatusCreated),
				createEntry(2, "2024-04-22", 0, task.StatusCreated),
			},
			wantResponse: []task.Task{
				createEntry(1, "2024-04-21", 0, task.StatusSnoozed),
				createEntry(2, "2024-04-22", 0, task.StatusSnoozed),
				createEntry(2, "2024-05-01", 0, task.StatusCreated),
				createEntry(1, "2024-05-01", 1, task.StatusCreated),
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

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/api/tasks/loadDay/%s", tc.today), nil)
			req = mux.SetURLVars(req, map[string]string{"date": tc.today})
			res := httptest.NewRecorder()

			h.ImportPastTasks(res, req)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			// Verify

			var resTasks []task.Task
			err := db.Model(task.Task{}).Find(&resTasks).Error
			if err != nil {
				t.Fatalf("Failed to fetch Tasks from DB for comparison: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(task.Task{}, "ID", "TaskID")); diff != "" {
				t.Errorf("ImportPastTasks() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func cleanupTaskDB(db *gorm.DB) *gorm.DB {
	return db.Exec("DELETE FROM tasks")
}
