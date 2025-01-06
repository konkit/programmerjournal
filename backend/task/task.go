package task

import "programmerjournal-backend/date"

type Status string

const (
	StatusCreated   Status = "Created"
	StatusDone             = "Done"
	StatusSnoozed          = "Snoozed"
	StatusMigrated         = "Migrated"
	StatusCancelled        = "Cancelled"
	StatusNote             = "Note"
)

type Task struct {
	ID           uint      `json:"id" sql:"AUTO_INCREMENT" gorm:"primarykey"`
	TaskID       string    `json:"taskID"`
	Title        string    `json:"title"`
	Status       Status    `json:"status"`
	CreatedDate  date.Date `json:"createdDate"`
	Update       string    `json:"update"`
	Description  string    `json:"description"`
	Rank         int       `json:"rank" sql:"DEFAULT:0"`
	SnoozedUntil date.Date `json:"snoozedUntil"`
}

func Clone(src Task) Task {
	newTask := Task{
		TaskID:      src.TaskID,
		Title:       src.Title,
		Status:      src.Status,
		CreatedDate: src.CreatedDate,
		Update:      src.Update,
		Description: src.Description,
	}
	return newTask
}

type TaskUpdate struct {
	Date   date.Date `json:"date"`
	Update string    `json:"update"`
	Status Status    `json:"status"`
}

type TaskSummary struct {
	Task    Task         `json:"task"`
	Updates []TaskUpdate `json:"updates"`
}

//type WeeklySummary struct {
//	TaskUpdates []TaskWeeklySummary `json:"taskUpdates"`
//}

type TaskWeeklySummary struct {
	Task    Task         `json:"task"`
	Updates []TaskUpdate `json:"updates"`
}
