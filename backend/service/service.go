package service

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type Task struct {
	ID           uint       `json:"id" sql:"AUTO_INCREMENT" gorm:"primarykey"`
	Title        string     `json:"title"`
	Done         bool       `json:"done"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	FinishedAt   *time.Time `json:"finished_at"`
	SnoozedUntil *time.Time `json:"snoozed_until"`
	Updates      []Update   `json:"updates" gorm:"foreignKey:TaskId"`
}

type Update struct {
	ID          uint      `json:"id" sql:"AUTO_INCREMENT" gorm:"primarykey"`
	TaskId      uint      `json:"task_id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	DoneToday   bool      `json:"done_today"`
}

type Service struct {
	Db *gorm.DB
}

func New(dbPath string) Service {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Task{})
	db.AutoMigrate(&Update{})

	return Service{
		Db: db,
	}
}

type CreateTaskEntry struct {
	Title string `json:"title"`
}

func (s *Service) CreateTask(newTaskEntry CreateTaskEntry) error {
	newTask := Task{
		Title: newTaskEntry.Title,
	}
	return s.Db.Create(&newTask).Error
}

type UpdateTitleEntry struct {
	Id    uint   `json:"id"`
	Title string `json:"title"`
}

func (s *Service) UpdateTaskTitle(entry UpdateTitleEntry) (*Task, error) {
	t := &Task{}
	result := s.Db.First(t, entry.Id)
	if result.Error != nil {
		return nil, fmt.Errorf("could not find entry with id %s", entry.Id)
	}

	t.Title = entry.Title

	s.Db.Save(t)

	return t, nil
}

type ListTaskResponse struct {
	Tasks []ListTaskEntry `json:"tasks"`
}

type ListTaskEntry struct {
	ID           uint              `json:"id"`
	Title        string            `json:"title"`
	Done         bool              `json:"done"`
	CreatedAt    string            `json:"created_at"`
	UpdatedAt    string            `json:"updated_at"`
	FinishedAt   string            `json:"finished_at"`
	SnoozedUntil string            `json:"snoozed_until"`
	Updates      []ListUpdateEntry `json:"updates"`
}

type ListUpdateEntry struct {
	Date        string `json:"date"`
	Description string `json:"description"`
	DoneToday   bool   `json:"done_today"`
}

func (s *Service) ListTasks(viewedDate string) (ListTaskResponse, error) {
	tasks := []Task{}
	err := s.Db.Model(Task{}).Preload("Updates").Find(&tasks).Error
	if err != nil {
		return ListTaskResponse{}, err
	}

	response := ListTaskResponse{
		Tasks: []ListTaskEntry{},
	}
	for _, task := range tasks {
		te := ListTaskEntry{
			ID:           task.ID,
			Title:        task.Title,
			Done:         task.Done,
			CreatedAt:    toDateString(&task.CreatedAt),
			UpdatedAt:    toDateString(&task.UpdatedAt),
			FinishedAt:   toDateString(task.FinishedAt),
			SnoozedUntil: toDateString(task.SnoozedUntil),
		}

		for _, update := range task.Updates {
			updateDate := toDateString(&update.Date)
			if updateDate == viewedDate {
				ue := ListUpdateEntry{
					Date:        updateDate,
					Description: update.Description,
					DoneToday:   update.DoneToday,
				}

				te.Updates = append(te.Updates, ue)
			}
		}

		response.Tasks = append(response.Tasks, te)
	}

	return response, nil
}

func toDateString(date *time.Time) string {
	if date == nil {
		return ""
	}
	return date.Format("2006-01-02")
}

func fromDateString(date string) (time.Time, error) {
	return time.Parse("2006-01-02", date)
}
