package task

import (
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"os"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/entry"
	"strconv"
	"testing"
)

func TestSetTaskDone(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := database.InitDB(dbTestPath)
	defer os.Remove(dbTestPath)
	es := database.NewEntryService(db)

	_, api := humatest.New(t)
	SetTaskDoneHandler(api, es)

	testCases := []struct {
		name         string
		initTask     entry.Entry
		wantResponse entry.Entry
		date         string
	}{
		{
			name: "set task done",
			initTask: entry.Entry{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      entry.StatusTaskCreated,
				CreatedDate: "2024-05-01",
				TaskUpdate:  "",
			},
			wantResponse: entry.Entry{
				TaskID:      "1234",
				Title:       "test 1",
				Status:      entry.StatusTaskDone,
				CreatedDate: "2024-05-01",
				TaskUpdate:  "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM entries")

			db.Create(&tc.initTask)

			insertedID := tc.initTask.ID

			e := struct{ Done bool }{
				Done: true,
			}
			url := fmt.Sprintf("/api/tasks/%d/setDone", insertedID)
			res := api.Patch(url, e)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status 201, got %d", res.Code)
			}

			var resTasks = entry.Entry{ID: insertedID}
			db.First(&resTasks)

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(entry.Entry{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestSetTaskDoneAndMoveToTheTop(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := database.InitDB(dbTestPath)
	defer os.Remove(dbTestPath)
	es := database.NewEntryService(db)

	_, api := humatest.New(t)
	SetTaskDoneHandler(api, es)

	createEntry := func(titleID int, rank int, status entry.Status) entry.Entry {
		return entry.Entry{
			TaskID:      strconv.Itoa(titleID),
			Title:       fmt.Sprintf("test %d", titleID),
			Status:      status,
			CreatedDate: "2024-05-01",
			TaskUpdate:  "",
			Rank:        rank,
		}
	}

	testCases := []struct {
		name              string
		initTasks         []entry.Entry
		modifiedTaskIndex uint
		wantResponses     []entry.Entry
		date              string
	}{
		{
			name: "one done modify last one",
			initTasks: []entry.Entry{
				createEntry(1, 0, entry.StatusTaskDone),
				createEntry(2, 1, entry.StatusTaskCreated),
				createEntry(3, 2, entry.StatusTaskCreated),
			},
			modifiedTaskIndex: 2,
			wantResponses: []entry.Entry{
				createEntry(1, 0, entry.StatusTaskDone),
				createEntry(3, 1, entry.StatusTaskDone),
				createEntry(2, 2, entry.StatusTaskCreated),
			},
		},
		{
			name: "no done modify first one",
			initTasks: []entry.Entry{
				createEntry(1, 0, entry.StatusTaskCreated),
				createEntry(2, 1, entry.StatusTaskCreated),
				createEntry(3, 2, entry.StatusTaskCreated),
			},
			modifiedTaskIndex: 0,
			wantResponses: []entry.Entry{
				createEntry(1, 0, entry.StatusTaskDone),
				createEntry(2, 1, entry.StatusTaskCreated),
				createEntry(3, 2, entry.StatusTaskCreated),
			},
		},
		{
			name: "no done modify last one",
			initTasks: []entry.Entry{
				createEntry(1, 0, entry.StatusTaskCreated),
				createEntry(2, 1, entry.StatusTaskCreated),
				createEntry(3, 2, entry.StatusTaskCreated),
			},
			modifiedTaskIndex: 2,
			wantResponses: []entry.Entry{
				createEntry(3, 0, entry.StatusTaskDone),
				createEntry(1, 1, entry.StatusTaskCreated),
				createEntry(2, 2, entry.StatusTaskCreated),
			},
		},
		{
			name: "one entry no done modify that one",
			initTasks: []entry.Entry{
				createEntry(1, 0, entry.StatusTaskCreated),
			},
			modifiedTaskIndex: 0,
			wantResponses: []entry.Entry{
				createEntry(1, 0, entry.StatusTaskDone),
			},
		},
		{
			name: "priority entry is not moved",
			initTasks: []entry.Entry{
				createEntry(-2, -2, entry.StatusTaskCreated),
				createEntry(0, 0, entry.StatusTaskCreated),
				createEntry(1, 1, entry.StatusTaskCreated),
			},
			modifiedTaskIndex: 0,
			wantResponses: []entry.Entry{
				createEntry(-2, -2, entry.StatusTaskDone),
				createEntry(0, 0, entry.StatusTaskCreated),
				createEntry(1, 1, entry.StatusTaskCreated),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM entries")

			var ids []uint
			for _, task := range tc.initTasks {
				db.Create(&task)
				ids = append(ids, task.ID)
			}

			e := struct{ Done bool }{
				Done: true,
			}
			url := fmt.Sprintf("/api/tasks/%d/setDone", ids[tc.modifiedTaskIndex])
			res := api.Patch(url, e)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status 201, got %d", res.Code)
			}

			resTasks, err := es.FindEntriesByDate("2024-05-01")
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponses, resTasks, cmpopts.IgnoreFields(entry.Entry{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
