package mcpserver

import (
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
)

type HistoryItem struct {
	EntryID uint            `json:"entryId"`
	Date    date.DateString `json:"date"`
	Status  entry.Status    `json:"status"`
	Update  string          `json:"update,omitempty"`
}

type TaskHistory struct {
	TaskID        string          `json:"taskId"`
	Title         string          `json:"title"`
	Description   string          `json:"description,omitempty"`
	CurrentStatus entry.Status    `json:"currentStatus"`
	DoneOn        date.DateString `json:"doneOn,omitempty" jsonschema:"date of the log where the task was marked done"`
	History       []HistoryItem   `json:"history"`
}

type ListDoneTasksInput struct {
	From string `json:"from" jsonschema:"first day of the range, YYYY-MM-DD"`
	To   string `json:"to" jsonschema:"last day of the range, YYYY-MM-DD"`
	Tag  string `json:"tag,omitempty" jsonschema:"only tasks with this tag, without the #"`
}

type ListDoneTasksOutput struct {
	Tasks []TaskHistory `json:"tasks"`
}

type GetTaskHistoryInput struct {
	TaskID  string `json:"taskId,omitempty" jsonschema:"taskId shared by all copies of the task"`
	EntryID uint   `json:"entryId,omitempty" jsonschema:"ID of any copy of the task; used when taskId is not given"`
}

type ListEntriesInput struct {
	Date string `json:"date" jsonschema:"log date: YYYY-MM-DD (day), YYYY-Www (week) or YYYY-MM (month)"`
}

type EntryItem struct {
	EntryID          uint            `json:"entryId"`
	TaskID           string          `json:"taskId,omitempty"`
	Date             date.DateString `json:"date"`
	Title            string          `json:"title"`
	Description      string          `json:"description,omitempty"`
	Status           entry.Status    `json:"status"`
	Update           string          `json:"update,omitempty"`
	TaskSnoozedUntil date.DateString `json:"taskSnoozedUntil,omitempty"`
}

type ListEntriesOutput struct {
	Entries []EntryItem `json:"entries"`
}

type SearchEntriesInput struct {
	Query  string `json:"query" jsonschema:"text to look for"`
	Status string `json:"status,omitempty" jsonschema:"only results in this current status: TaskCreated, TaskDone, TaskSnoozed, TaskMigrated, TaskCancelled or Note"`
	Limit  int    `json:"limit,omitempty" jsonschema:"maximum number of results, most recent first; default 50"`
}

type SearchEntriesOutput struct {
	Results []EntryItem `json:"results"`
}

type ListTagsOutput struct {
	Tags []string `json:"tags"`
}
