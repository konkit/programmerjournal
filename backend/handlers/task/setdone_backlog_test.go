package task

import (
	"fmt"
	"github.com/danielgtaylor/huma/v2/humatest"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/entry"
	"testing"
)

func TestSetTaskDoneUpdatesBacklogs(t *testing.T) {
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	es := database.NewEntryService(db)

	_, api := humatest.New(t)
	SetTaskDoneHandler(api, es)

	taskID := "linked-task-123"

	// 1. Create a weekly entry for the same week as the daily task
	// 2024-05-06 is Monday of week 19
	weeklyEntry := entry.Entry{
		TaskID:      taskID,
		Title:       "Weekly Task",
		Status:      entry.StatusTaskSnoozed,
		CreatedDate: "2024-W19",
	}
	db.Create(&weeklyEntry)

	// 2. Create a monthly entry for the same month
	monthlyEntry := entry.Entry{
		TaskID:      taskID,
		Title:       "Monthly Task",
		Status:      entry.StatusTaskSnoozed,
		CreatedDate: "2024-05",
	}
	db.Create(&monthlyEntry)

	// 3. Create a daily entry
	dailyEntry := entry.Entry{
		TaskID:      taskID,
		Title:       "Daily Task",
		Status:      entry.StatusTaskCreated,
		CreatedDate: "2024-05-06",
	}
	db.Create(&dailyEntry)

	// 4. Create an entry from a DIFFERENT week to ensure it's NOT updated
	otherWeeklyEntry := entry.Entry{
		TaskID:      taskID,
		Title:       "Past Weekly Task",
		Status:      entry.StatusTaskSnoozed,
		CreatedDate: "2024-W18",
	}
	db.Create(&otherWeeklyEntry)

	// Mark daily as done
	url := fmt.Sprintf("/api/tasks/%d/setDone", dailyEntry.ID)
	res := api.Patch(url, struct{ Done bool }{Done: true})

	if res.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", res.Code)
	}

	// Verify all entries
	var checkedDaily entry.Entry
	db.First(&checkedDaily, dailyEntry.ID)
	if checkedDaily.Status != entry.StatusTaskDone {
		t.Errorf("Daily task status expected Done, got %s", checkedDaily.Status)
	}

	var checkedWeekly entry.Entry
	db.First(&checkedWeekly, weeklyEntry.ID)
	if checkedWeekly.Status != entry.StatusTaskDone {
		t.Errorf("Weekly task status expected Done, got %s", checkedWeekly.Status)
	}

	var checkedMonthly entry.Entry
	db.First(&checkedMonthly, monthlyEntry.ID)
	if checkedMonthly.Status != entry.StatusTaskDone {
		t.Errorf("Monthly task status expected Done, got %s", checkedMonthly.Status)
	}

	var checkedOtherWeekly entry.Entry
	db.First(&checkedOtherWeekly, otherWeeklyEntry.ID)
	if checkedOtherWeekly.Status != entry.StatusTaskSnoozed {
		t.Errorf("Other weekly task status expected Snoozed (unmodified), got %s", checkedOtherWeekly.Status)
	}
}
