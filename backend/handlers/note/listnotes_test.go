package note

import (
	"encoding/json"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/entry"
	"testing"

	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestListNotes(t *testing.T) {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	es := database.NewEntryService(db)

	_, api := humatest.New(t)
	ListNotesHandler(api, es)

	testCases := []struct {
		name         string
		initTasks    []entry.Entry
		wantResponse []entry.Entry
	}{
		{
			name: "list notes",
			initTasks: []entry.Entry{
				{
					Title:       "note 1",
					Status:      entry.StatusNote,
					CreatedDate: "2024-05-01",
					TaskUpdate:  "",
					Rank:        0,
				},
				{
					Title:       "task 1",
					Status:      entry.StatusTaskCreated,
					CreatedDate: "2024-05-01",
					TaskUpdate:  "",
					Rank:        0,
				},
				{
					Title:       "note 2",
					Status:      entry.StatusNote,
					CreatedDate: "2024-05-02",
					TaskUpdate:  "",
					Rank:        0,
				},
			},
			wantResponse: []entry.Entry{
				{
					Title:       "note 2",
					Status:      entry.StatusNote,
					CreatedDate: "2024-05-02",
					TaskUpdate:  "",
					Rank:        0,
				},
				{
					Title:       "note 1",
					Status:      entry.StatusNote,
					CreatedDate: "2024-05-01",
					TaskUpdate:  "",
					Rank:        0,
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM entries")

			for _, task := range tc.initTasks {
				db.Create(&task)
			}

			url := "/api/notes"
			res := api.Get(url)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status 200, got %d", res.Code)
			}

			var responseBody ListNotesResponse
			err := json.NewDecoder(res.Body).Decode(&responseBody)
			if err != nil {
				t.Fatalf("Failed to deserialize response: %v", err)
			}

			if diff := cmp.Diff(tc.wantResponse, responseBody.Body.Notes, cmpopts.IgnoreFields(entry.Entry{}, "ID", "TaskID")); diff != "" {
				t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
