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
	Rank        int    `json:"rank" sql:"DEFAULT:0"`
}

func Clone(src Task) Task {
	newTask := Task{
		Title:       src.Title,
		TaskID:      src.TaskID,
		Status:      src.Status,
		Update:      src.Update,
		CreatedDate: src.CreatedDate,
	}
	return newTask
}
