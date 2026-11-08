package mcpserver

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// A task created on May 1st, snoozed for two days with progress notes and done on May 3rd.
// Its weekly copy was marked done along with the daily one.
var progressedTask = []entry.Entry{
	{TaskID: "t1", Title: "Write report #work", Status: entry.StatusTaskSnoozed, CreatedDate: "2024-05-01", TaskUpdate: "outline"},
	{TaskID: "t1", Title: "Write report #work", Status: entry.StatusTaskSnoozed, CreatedDate: "2024-05-02", TaskUpdate: "first draft"},
	{TaskID: "t1", Title: "Write report #work", Status: entry.StatusTaskDone, CreatedDate: "2024-05-03", TaskUpdate: "sent"},
	{TaskID: "t1", Title: "Write report #work", Status: entry.StatusTaskDone, CreatedDate: "2024-W18"},
}

func setup(t *testing.T, entries ...[]entry.Entry) *database.EntryService {
	t.Helper()
	db, err := database.InitDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	es := database.NewEntryService(db)
	for _, group := range entries {
		for _, e := range group {
			if err := es.InsertEntry(&e); err != nil {
				t.Fatalf("Failed to insert entry: %v", err)
			}
		}
	}
	return es
}

func historyDates(th TaskHistory) []date.DateString {
	var dates []date.DateString
	for _, h := range th.History {
		dates = append(dates, h.Date)
	}
	return dates
}

func TestListDoneTasks(t *testing.T) {
	es := setup(t, progressedTask, []entry.Entry{
		{TaskID: "t2", Title: "Plan sprint", Status: entry.StatusTaskDone, CreatedDate: "2024-05-20"},
		{TaskID: "t3", Title: "Monthly review", Status: entry.StatusTaskDone, CreatedDate: "2024-04"},
		{TaskID: "t4", Title: "Still open #work", Status: entry.StatusTaskCreated, CreatedDate: "2024-05-02"},
	})

	testCases := []struct {
		name      string
		input     ListDoneTasksInput
		wantTasks []string
	}{
		{name: "range with the task counts it once", input: ListDoneTasksInput{From: "2024-05-01", To: "2024-05-10"}, wantTasks: []string{"t1"}},
		{name: "range covering only the weekly copy", input: ListDoneTasksInput{From: "2024-05-04", To: "2024-05-05"}, wantTasks: nil},
		{name: "monthly task overlapping the range", input: ListDoneTasksInput{From: "2024-04-30", To: "2024-05-01"}, wantTasks: []string{"t3"}},
		{name: "ordered by completion", input: ListDoneTasksInput{From: "2024-04-01", To: "2024-05-31"}, wantTasks: []string{"t3", "t1", "t2"}},
		{name: "filtered by tag", input: ListDoneTasksInput{From: "2024-04-01", To: "2024-05-31", Tag: "#Work"}, wantTasks: []string{"t1"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := ListDoneTasks(es, tc.input)
			if err != nil {
				t.Fatalf("ListDoneTasks: %v", err)
			}
			var got []string
			for _, task := range out.Tasks {
				got = append(got, task.TaskID)
			}
			if diff := cmp.Diff(tc.wantTasks, got); diff != "" {
				t.Errorf("tasks mismatch (-want +got):\n%s", diff)
			}
		})
	}

	out, _ := ListDoneTasks(es, ListDoneTasksInput{From: "2024-05-03", To: "2024-05-03"})
	task := out.Tasks[0]
	if task.DoneOn != "2024-05-03" || task.CurrentStatus != entry.StatusTaskDone {
		t.Errorf("doneOn %q, status %q; want 2024-05-03, TaskDone", task.DoneOn, task.CurrentStatus)
	}
	// 2024-W18 starts on Monday, April 29th, so the weekly copy comes first
	wantDates := []date.DateString{"2024-W18", "2024-05-01", "2024-05-02", "2024-05-03"}
	if diff := cmp.Diff(wantDates, historyDates(task)); diff != "" {
		t.Errorf("history mismatch (-want +got):\n%s", diff)
	}

	if _, err := ListDoneTasks(es, ListDoneTasksInput{From: "2024-W18", To: "2024-05-03"}); err == nil {
		t.Error("expected an error for a non-day from date")
	}
}

func TestSortChain(t *testing.T) {
	chain := []entry.Entry{
		{CreatedDate: "2024-05-02"}, {CreatedDate: "2024-06"}, {CreatedDate: "2024-W18"},
		{CreatedDate: "2024-05"}, {CreatedDate: "2024-04-30"},
	}
	sortChain(chain)

	var got []date.DateString
	for _, e := range chain {
		got = append(got, e.CreatedDate)
	}
	// 2024-W18 starts on Monday, April 29th
	want := []date.DateString{"2024-W18", "2024-04-30", "2024-05", "2024-05-02", "2024-06"}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("order mismatch (-want +got):\n%s", diff)
	}
}

func TestSearchEntries(t *testing.T) {
	es := setup(t, progressedTask, []entry.Entry{
		{Title: "Draft ideas", Status: entry.StatusNote, CreatedDate: "2024-05-04"},
		{TaskID: "t5", Title: "Unrelated", Status: entry.StatusTaskCreated, CreatedDate: "2024-05-04"},
	})

	out, err := SearchEntries(es, SearchEntriesInput{Query: "draft"})
	if err != nil {
		t.Fatalf("SearchEntries: %v", err)
	}
	// The task matched on an old update, but is shown in its latest state
	want := []string{"Draft ideas/Note/2024-05-04", "Write report #work/TaskDone/2024-05-03"}
	var got []string
	for _, r := range out.Results {
		got = append(got, r.Title+"/"+string(r.Status)+"/"+string(r.Date))
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("results mismatch (-want +got):\n%s", diff)
	}

	out, _ = SearchEntries(es, SearchEntriesInput{Query: "draft", Status: entry.StatusNote})
	if len(out.Results) != 1 || out.Results[0].Status != entry.StatusNote {
		t.Errorf("status filter: got %+v", out.Results)
	}
}

func TestOverHTTP(t *testing.T) {
	es := setup(t, progressedTask)
	srv := httptest.NewServer(NewHandler(es))
	defer srv.Close()

	ctx := context.Background()
	client := mcp.NewClient(&mcp.Implementation{Name: "test", Version: "1.0.0"}, nil)
	session, err := client.Connect(ctx, &mcp.StreamableClientTransport{Endpoint: srv.URL}, nil)
	if err != nil {
		t.Fatalf("Connect: %v", err)
	}
	defer session.Close()

	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("ListTools: %v", err)
	}
	var names []string
	for _, tool := range tools.Tools {
		names = append(names, tool.Name)
		if tool.Annotations == nil || !tool.Annotations.ReadOnlyHint {
			t.Errorf("tool %s is not marked read-only", tool.Name)
		}
	}
	wantNames := []string{"get_task_history", "list_done_tasks", "list_entries", "list_tags", "search_entries"}
	if diff := cmp.Diff(wantNames, names); diff != "" {
		t.Errorf("tools mismatch (-want +got):\n%s", diff)
	}

	res, err := session.CallTool(ctx, &mcp.CallToolParams{Name: "get_task_history", Arguments: map[string]any{"taskId": "t1"}})
	if err != nil || res.IsError {
		t.Fatalf("CallTool get_task_history: %v %+v", err, res)
	}
	var th TaskHistory
	raw, _ := json.Marshal(res.StructuredContent)
	if err := json.Unmarshal(raw, &th); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if th.DoneOn != "2024-05-03" || len(th.History) != 4 || th.History[2].Update != "first draft" {
		t.Errorf("unexpected history: %+v", th)
	}

	res, err = session.CallTool(ctx, &mcp.CallToolParams{Name: "list_entries", Arguments: map[string]any{"date": "yesterday"}})
	if err != nil || !res.IsError {
		t.Errorf("expected a tool error for a bad date, got %v %+v", err, res)
	}
}
