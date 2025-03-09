package tags

import (
	"encoding/json"
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"os"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/entry"
	"testing"
)

func TestGetTagOld(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := database.InitDB(dbTestPath)
	defer os.Remove(dbTestPath)
	es := database.NewEntryService(db)

	_, api := humatest.New(t)
	GetTagsHandler(api, db)

	db.Exec("DELETE FROM entries")

	initTasks := []entry.Entry{
		{
			ID:     0,
			Title:  "test task #taskone",
			Status: entry.StatusTaskCreated,
		},
	}

	for _, task := range initTasks {
		err := es.InsertEntry(&task)
		if err != nil {
			t.Fatalf("Error inserting entry: %v", err)
		}
	}

	url := fmt.Sprintf("/api/tags/%s", "taskone")
	res := api.Get(url)

	if res.Code != http.StatusOK {
		t.Fatalf("Expected status OK, got %d", res.Code)
	}

	wantResponse := []entry.Entry{
		{
			ID:     0,
			Title:  "test task #taskone",
			Status: entry.StatusTaskCreated,
		},
	}
	resTasks := []entry.Entry{}
	err := json.NewDecoder(res.Body).Decode(&resTasks)
	if err != nil {
		t.Fatalf("Failed to deserialize response: %v", err)
	}

	if diff := cmp.Diff(wantResponse, resTasks, cmpopts.IgnoreFields(entry.Entry{}, "ID", "TaskID")); diff != "" {
		t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
	}
}

func TestGetTag(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := database.InitDB(dbTestPath)
	defer os.Remove(dbTestPath)
	es := database.NewEntryService(db)

	_, api := humatest.New(t)
	GetTagsHandler(api, db)

	testCases := []struct {
		name         string
		initTasks    []entry.Entry
		wantTag      string
		wantResponse []entry.Entry
	}{
		{
			name: "one tag",
			initTasks: []entry.Entry{
				{
					ID:     0,
					Title:  "test task #taskone",
					Status: entry.StatusTaskCreated,
				},
			},
			wantTag: "taskone",
			wantResponse: []entry.Entry{
				{
					ID:     0,
					Title:  "test task #taskone",
					Status: entry.StatusTaskCreated,
				},
			},
		},
		{
			name: "one of two tags",
			initTasks: []entry.Entry{
				{
					ID:     0,
					Title:  "test task #taskone",
					Status: entry.StatusTaskCreated,
					Rank:   0,
				},
				{
					ID:     1,
					Title:  "test task #tasktwo",
					Status: entry.StatusTaskCreated,
					Rank:   1,
				},
			},
			wantTag: "tasktwo",
			wantResponse: []entry.Entry{
				{
					ID:     1,
					Title:  "test task #tasktwo",
					Status: entry.StatusTaskCreated,
					Rank:   1,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM entries")

			for _, task := range tc.initTasks {
				err := es.InsertEntry(&task)
				if err != nil {
					t.Fatalf("Error inserting entry: %v", err)
				}
			}

			url := fmt.Sprintf("/api/tags/%s", tc.wantTag)
			res := api.Get(url)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			resTasks := []entry.Entry{}
			err := json.NewDecoder(res.Body).Decode(&resTasks)
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, resTasks, cmpopts.IgnoreFields(entry.Entry{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
