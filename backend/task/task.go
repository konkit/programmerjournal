package task

type Status string

const (
	StatusCreated   Status = "Created"
	StatusDone             = "Done"
	StatusSnoozed          = "Snoozed"
	StatusCancelled        = "Cancelled"
)

type Task struct {
	ID          uint   `json:"id" sql:"AUTO_INCREMENT" gorm:"primarykey"`
	TaskID      string `json:"taskID"`
	Title       string `json:"title"`
	Status      Status `json:"status"`
	CreatedDate string `json:"createdDate"`
	Update      string `json:"update"`
	Description string `json:"description"`
	Rank        int    `json:"rank" sql:"DEFAULT:0"`
}

func Clone(src Task) Task {
	newTask := Task{
		Title:       src.Title,
		TaskID:      src.TaskID,
		Status:      src.Status,
		Update:      src.Update,
		Description: src.Description,
		CreatedDate: src.CreatedDate,
	}
	return newTask
}

type TaskUpdate struct {
	Date   string `json:"date"`
	Update string `json:"update"`
}

type TaskSummary struct {
	Task    Task         `json:"task"`
	Updates []TaskUpdate `json:"updates"`
}
