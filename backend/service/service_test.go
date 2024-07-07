package service_test

import (
	"github.com/google/go-cmp/cmp"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"programmerjournal-backend/service"
	"testing"
	"time"
)

func TestListTasks(t *testing.T) {
	db := initTestDb()
	s := service.Service{Db: db}

	testCases := []struct {
		name         string
		initTasks    []service.Task
		viewedDate   string
		wantResponse service.ListTaskResponse
	}{
		{
			name: "list single task",
			initTasks: []service.Task{
				{
					Title:      "test 1",
					Done:       false,
					FinishedAt: newDatePtr(2024, time.June, 1),
				},
			},
			viewedDate: "2024-06-01",
			wantResponse: service.ListTaskResponse{
				Tasks: []service.ListTaskEntry{
					{
						Title:      "test 1",
						Done:       false,
						CreatedAt:  time.Now().Format("2006-01-02"),
						UpdatedAt:  time.Now().Format("2006-01-02"),
						FinishedAt: "2024-06-01",
					},
				},
			},
		},
		{
			name: "list single task with a single updates",
			initTasks: []service.Task{
				{
					Title:      "test 1",
					Done:       false,
					CreatedAt:  newDate(2024, time.June, 1),
					UpdatedAt:  newDate(2024, time.June, 1),
					FinishedAt: newDatePtr(2024, time.June, 1),
					Updates: []service.Update{
						{
							Date:        newDate(2024, time.June, 1),
							Description: "description of the task",
							DoneToday:   false,
						},
					},
				},
			},
			viewedDate: "2024-06-01",
			wantResponse: service.ListTaskResponse{
				Tasks: []service.ListTaskEntry{
					{
						Title:      "test 1",
						Done:       false,
						CreatedAt:  "2024-06-01",
						UpdatedAt:  "2024-06-01",
						FinishedAt: "2024-06-01",
						Updates: []service.ListUpdateEntry{
							{
								Date:        "2024-06-01",
								Description: "description of the task",
								DoneToday:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "list single task with updates from different dates",
			initTasks: []service.Task{
				{
					Title:      "test 1",
					Done:       false,
					CreatedAt:  newDate(2024, time.June, 1),
					UpdatedAt:  newDate(2024, time.June, 1),
					FinishedAt: newDatePtr(2024, time.June, 1),
					Updates: []service.Update{
						{
							Date:        newDate(2024, time.May, 31),
							Description: "description of the task",
							DoneToday:   false,
						},
						{
							Date:        newDate(2024, time.June, 1),
							Description: "description of the task",
							DoneToday:   false,
						},
					},
				},
			},
			viewedDate: "2024-06-01",
			wantResponse: service.ListTaskResponse{
				Tasks: []service.ListTaskEntry{
					{
						Title:      "test 1",
						Done:       false,
						CreatedAt:  "2024-06-01",
						UpdatedAt:  "2024-06-01",
						FinishedAt: "2024-06-01",
						Updates: []service.ListUpdateEntry{
							{
								Date:        "2024-06-01",
								Description: "description of the task",
								DoneToday:   false,
							},
						},
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, task := range tc.initTasks {
				s.Db.Create(&task)
			}

			response, err := s.ListTasks(tc.viewedDate)
			if err != nil {
				t.Errorf("ListTasks() failed: %v", err)
			}

			if len(response.Tasks) != len(tc.wantResponse.Tasks) {
				t.Errorf("got %d tasks, want %d", len(response.Tasks), len(tc.wantResponse.Tasks))
			}

			if diff := cmp.Diff(tc.wantResponse, response); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}

			db.Exec("DELETE FROM tasks")
		})
	}
}

func initTestDb() *gorm.DB {
	dbPath := "../test.db"
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&service.Task{})
	db.AutoMigrate(&service.Update{})

	return db
}

func newDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func newDatePtr(year int, month time.Month, day int) *time.Time {
	date := newDate(year, month, day)
	return &date
}
