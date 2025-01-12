package entry

import "programmerjournal-backend/model/date"

type Status string

const (
	StatusTaskCreated   Status = "TaskCreated"
	StatusTaskDone             = "TaskDone"
	StatusTaskSnoozed          = "TaskSnoozed"
	StatusTaskMigrated         = "TaskMigrated"
	StatusTaskCancelled        = "TaskCancelled"
	StatusNote                 = "Note"
)

type Entry struct {
	ID               uint            `json:"id" sql:"AUTO_INCREMENT" gorm:"primarykey"`
	Title            string          `json:"title"`
	Status           Status          `json:"status"`
	CreatedDate      date.DateString `json:"createdDate"`
	Description      string          `json:"description"`
	Rank             int             `json:"rank" sql:"DEFAULT:0"`
	TaskID           string          `json:"taskID"`
	TaskUpdate       string          `json:"taskUpdate"`
	TaskSnoozedUntil date.DateString `json:"taskSnoozedUntil"`
}

func Clone(src Entry) Entry {
	newEntry := Entry{
		Title:       src.Title,
		Status:      src.Status,
		CreatedDate: src.CreatedDate,
		Description: src.Description,
		TaskID:      src.TaskID,
		TaskUpdate:  src.TaskUpdate,
	}
	return newEntry
}

type TaskUpdate struct {
	Date   date.DateString `json:"date"`
	Update string          `json:"update"`
	Status Status          `json:"status"`
}

type TaskSummary struct {
	TaskEntry Entry        `json:"taskEntry"`
	Updates   []TaskUpdate `json:"updates"`
}
