// Package mcpserver exposes the journal as a read-only MCP server, mounted at /mcp.
package mcpserver

import (
	"context"
	"fmt"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
	"sort"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const instructions = `Programmer journal: a bullet journal with daily (YYYY-MM-DD), weekly (YYYY-Www) and monthly (YYYY-MM) logs.
A task that is snoozed or migrated is copied to a new date; the copies share a taskID and each copy can carry a short progress update.
The chain of copies, ordered by date, is the task's history. A task's current status is the status of its last copy.
Tags are #hashtags in task titles.`

func NewHandler(es *database.EntryService) http.Handler {
	server := NewServer(es)
	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server { return server }, &mcp.StreamableHTTPOptions{
		Stateless:    true,
		JSONResponse: true,
	})
}

func NewServer(es *database.EntryService) *mcp.Server {
	server := mcp.NewServer(&mcp.Implementation{Name: "programmer-journal", Version: "1.0.0"}, &mcp.ServerOptions{
		Instructions: instructions,
	})

	readOnly := &mcp.ToolAnnotations{ReadOnlyHint: true}

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_done_tasks",
		Description: "List tasks completed within a date range (inclusive), each with its full progress history.",
		Annotations: readOnly,
	}, func(ctx context.Context, req *mcp.CallToolRequest, in ListDoneTasksInput) (*mcp.CallToolResult, ListDoneTasksOutput, error) {
		out, err := ListDoneTasks(es, in)
		return nil, out, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_task_history",
		Description: "Get the progress history of one task: every dated copy of it with its status and update.",
		Annotations: readOnly,
	}, func(ctx context.Context, req *mcp.CallToolRequest, in GetTaskHistoryInput) (*mcp.CallToolResult, TaskHistory, error) {
		out, err := GetTaskHistory(es, in)
		return nil, out, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_entries",
		Description: "List all tasks and notes of one daily, weekly or monthly log, including snoozed and migrated copies.",
		Annotations: readOnly,
	}, func(ctx context.Context, req *mcp.CallToolRequest, in ListEntriesInput) (*mcp.CallToolResult, ListEntriesOutput, error) {
		out, err := ListEntries(es, in)
		return nil, out, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "search_entries",
		Description: "Search tasks and notes by text in their title, description or progress updates. Returns one result per task, in its current state.",
		Annotations: readOnly,
	}, func(ctx context.Context, req *mcp.CallToolRequest, in SearchEntriesInput) (*mcp.CallToolResult, SearchEntriesOutput, error) {
		out, err := SearchEntries(es, in)
		return nil, out, err
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_tags",
		Description: "List all tags used in the journal.",
		Annotations: readOnly,
	}, func(ctx context.Context, req *mcp.CallToolRequest, in struct{}) (*mcp.CallToolResult, ListTagsOutput, error) {
		names, err := es.ListTagNames()
		return nil, ListTagsOutput{Tags: names}, err
	})

	return server
}

// ListDoneTasks returns tasks whose completion falls in the range. A task done in a weekly or
// monthly log counts when that week or month overlaps the range.
func ListDoneTasks(es *database.EntryService, in ListDoneTasksInput) (ListDoneTasksOutput, error) {
	out := ListDoneTasksOutput{Tasks: []TaskHistory{}}

	from, err := date.ParseDayDate(date.DateString(in.From))
	if err != nil {
		return out, fmt.Errorf("from: %w", err)
	}
	to, err := date.ParseDayDate(date.DateString(in.To))
	if err != nil {
		return out, fmt.Errorf("to: %w", err)
	}

	doneEntries, err := es.FindEntriesByStatus(entry.StatusTaskDone)
	if err != nil {
		return out, err
	}

	if in.Tag != "" {
		taggedIDs, err := es.FindEntryIDsByTag(strings.ToLower(strings.TrimPrefix(in.Tag, "#")))
		if err != nil {
			return out, err
		}
		doneEntries = filter(doneEntries, func(e entry.Entry) bool { return taggedIDs[e.ID] })
	}

	// Marking a daily task done also marks its weekly and monthly copies done, so group by task
	for _, done := range groupByTask(doneEntries) {
		completion := completionEntry(done)
		start, end, err := date.GetDayRange(completion.CreatedDate)
		if err != nil {
			continue
		}
		if start.Value > to.Value || end.Value < from.Value {
			continue
		}

		chain, err := taskChain(es, completion)
		if err != nil {
			return out, err
		}
		history := buildHistory(chain)
		history.DoneOn = completion.CreatedDate
		out.Tasks = append(out.Tasks, history)
	}

	sort.SliceStable(out.Tasks, func(i, j int) bool {
		return lessDate(out.Tasks[i].DoneOn, out.Tasks[j].DoneOn)
	})
	return out, nil
}

func GetTaskHistory(es *database.EntryService, in GetTaskHistoryInput) (TaskHistory, error) {
	var chain []entry.Entry
	switch {
	case in.TaskID != "":
		tasks, err := es.FindTasksByTaskID(in.TaskID)
		if err != nil {
			return TaskHistory{}, err
		}
		chain = tasks
	case in.EntryID != 0:
		e, err := es.GetEntryByID(in.EntryID)
		if err != nil {
			return TaskHistory{}, fmt.Errorf("entry %d: %w", in.EntryID, err)
		}
		chain, err = taskChain(es, e)
		if err != nil {
			return TaskHistory{}, err
		}
	default:
		return TaskHistory{}, fmt.Errorf("either taskId or entryId is required")
	}

	if len(chain) == 0 {
		return TaskHistory{}, fmt.Errorf("task %q not found", in.TaskID)
	}
	return buildHistory(chain), nil
}

func ListEntries(es *database.EntryService, in ListEntriesInput) (ListEntriesOutput, error) {
	out := ListEntriesOutput{Entries: []EntryItem{}}
	if date.GetDateType(in.Date) == date.DateTypeUnrecognized {
		return out, fmt.Errorf("unrecognized date format: %v", in.Date)
	}

	entries, err := es.FindEntriesByDate(date.DateString(in.Date))
	if err != nil {
		return out, err
	}
	for _, e := range entries {
		out.Entries = append(out.Entries, toEntryItem(e))
	}
	return out, nil
}

func SearchEntries(es *database.EntryService, in SearchEntriesInput) (SearchEntriesOutput, error) {
	out := SearchEntriesOutput{Results: []EntryItem{}}
	if strings.TrimSpace(in.Query) == "" {
		return out, fmt.Errorf("query is required")
	}
	limit := in.Limit
	if limit <= 0 {
		limit = 50
	}

	matches, err := es.SearchEntries(in.Query)
	if err != nil {
		return out, err
	}

	// A match in any copy of a task stands for the whole task, shown as its latest copy
	for _, group := range groupByTask(matches) {
		latest := group[len(group)-1]
		if latest.TaskID != "" {
			chain, err := es.FindTasksByTaskID(latest.TaskID)
			if err != nil {
				return out, err
			}
			sortChain(chain)
			latest = chain[len(chain)-1]
		}
		if in.Status != "" && string(latest.Status) != in.Status {
			continue
		}
		out.Results = append(out.Results, toEntryItem(latest))
	}

	sort.SliceStable(out.Results, func(i, j int) bool {
		return lessDate(out.Results[j].Date, out.Results[i].Date)
	})
	if len(out.Results) > limit {
		out.Results = out.Results[:limit]
	}
	return out, nil
}

func toEntryItem(e entry.Entry) EntryItem {
	return EntryItem{
		EntryID:          e.ID,
		TaskID:           e.TaskID,
		Date:             e.CreatedDate,
		Title:            e.Title,
		Description:      e.Description,
		Status:           e.Status,
		Update:           e.TaskUpdate,
		TaskSnoozedUntil: e.TaskSnoozedUntil,
	}
}

// taskChain returns all copies of the task the entry belongs to, ordered by date.
func taskChain(es *database.EntryService, e entry.Entry) ([]entry.Entry, error) {
	if e.TaskID == "" {
		return []entry.Entry{e}, nil
	}
	chain, err := es.FindTasksByTaskID(e.TaskID)
	if err != nil {
		return nil, err
	}
	sortChain(chain)
	return chain, nil
}

func buildHistory(chain []entry.Entry) TaskHistory {
	sortChain(chain)
	latest := chain[len(chain)-1]
	th := TaskHistory{
		TaskID:        latest.TaskID,
		Title:         latest.Title,
		Description:   latest.Description,
		CurrentStatus: latest.Status,
		History:       []HistoryItem{},
	}
	for _, e := range chain {
		th.History = append(th.History, HistoryItem{
			EntryID: e.ID,
			Date:    e.CreatedDate,
			Status:  e.Status,
			Update:  e.TaskUpdate,
		})
	}
	if latest.Status == entry.StatusTaskDone {
		th.DoneOn = completionEntry(chain).CreatedDate
	}
	return th
}

// completionEntry picks the copy where the task was actually ticked: the latest done daily copy
// if there is one, otherwise the latest done weekly or monthly copy.
func completionEntry(entries []entry.Entry) entry.Entry {
	sortChain(entries)
	var best *entry.Entry
	for i := range entries {
		e := &entries[i]
		if e.Status != entry.StatusTaskDone {
			continue
		}
		if best == nil || date.GetDateType(string(e.CreatedDate)) == date.DateTypeDay ||
			date.GetDateType(string(best.CreatedDate)) != date.DateTypeDay {
			best = e
		}
	}
	if best == nil {
		return entries[len(entries)-1]
	}
	return *best
}

// groupByTask groups entries by taskID, keeping the order in which tasks first appear.
// Entries without a taskID (notes) are groups of their own. Each group is ordered by date.
func groupByTask(entries []entry.Entry) [][]entry.Entry {
	var keys []string
	groups := map[string][]entry.Entry{}
	for _, e := range entries {
		key := e.TaskID
		if key == "" {
			key = fmt.Sprintf("entry:%d", e.ID)
		}
		if _, ok := groups[key]; !ok {
			keys = append(keys, key)
		}
		groups[key] = append(groups[key], e)
	}

	result := make([][]entry.Entry, 0, len(keys))
	for _, key := range keys {
		sortChain(groups[key])
		result = append(result, groups[key])
	}
	return result
}

// sortChain orders entries by the first day of their log; a month comes before its weeks and a
// week before its days.
func sortChain(entries []entry.Entry) {
	sort.SliceStable(entries, func(i, j int) bool {
		return lessDate(entries[i].CreatedDate, entries[j].CreatedDate)
	})
}

func lessDate(a, b date.DateString) bool {
	aStart, _, aErr := date.GetDayRange(a)
	bStart, _, bErr := date.GetDayRange(b)
	if aErr != nil || bErr != nil || aStart.Value == bStart.Value {
		return granularity(a) < granularity(b)
	}
	return aStart.Value < bStart.Value
}

func granularity(d date.DateString) int {
	switch date.GetDateType(string(d)) {
	case date.DateTypeMonth:
		return 0
	case date.DateTypeWeek:
		return 1
	default:
		return 2
	}
}

func filter(entries []entry.Entry, keep func(entry.Entry) bool) []entry.Entry {
	var result []entry.Entry
	for _, e := range entries {
		if keep(e) {
			result = append(result, e)
		}
	}
	return result
}
