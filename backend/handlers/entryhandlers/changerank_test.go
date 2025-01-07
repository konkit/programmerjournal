package entryhandlers

import (
	"encoding/json"
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"net/http"
	"os"
	"programmerjournal-backend/model/entry"
	"strconv"
	"testing"
)

func TestChangeRank(t *testing.T) {
	dbTestPath := "./test.db"
	db, _ := entry.InitDB(dbTestPath)
	defer os.Remove(dbTestPath)

	dbRepo, _ := entry.NewRepository(db)

	_, api := humatest.New(t)
	ChangeRank(api, dbRepo)
	ListEntries(api, dbRepo)

	createEntry := func(titleID int, rank int) entry.Entry {
		return entry.Entry{
			TaskID:      strconv.Itoa(titleID),
			Title:       fmt.Sprintf("test %d", titleID),
			Status:      entry.StatusTaskCreated,
			CreatedDate: "2024-05-01",
			TaskUpdate:  "",
			Rank:        rank,
		}
	}

	testCases := []struct {
		name         string
		initTasks    []entry.Entry
		oldRank      int
		newRank      int
		wantResponse []entry.Entry
	}{
		{
			name: "change rank in the middle",
			initTasks: []entry.Entry{
				createEntry(0, 0),
				createEntry(1, 1),
				createEntry(2, 2),
				createEntry(3, 3),
			},
			oldRank: 1,
			newRank: 3,
			wantResponse: []entry.Entry{
				createEntry(0, 0),
				createEntry(2, 1),
				createEntry(3, 2),
				createEntry(1, 3),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db.Exec("DELETE FROM tasks")

			var createdTasks []entry.Entry
			for _, tt := range tc.initTasks {
				db.Create(&tt)
				createdTasks = append(createdTasks, tt)
			}

			changedID := createdTasks[tc.oldRank].ID
			e := struct{ NewRank int }{
				NewRank: tc.newRank,
			}
			url := fmt.Sprintf("/api/entries/%d/changeRank", changedID)
			res := api.Patch(url, e)

			if res.Code != http.StatusOK {
				t.Fatalf("Expected status OK, got %d", res.Code)
			}

			listUrl := fmt.Sprintf("/api/entries/list/%s", "2024-05-01")
			res = api.Get(listUrl)

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
