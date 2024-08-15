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

type ListTaskResponse struct {
	Tasks []ListTaskEntry `json:"tasks"`
}

type ListTaskEntry struct {
	ID           uint   `json:"id"`
	Title        string `json:"title"`
	Done         bool   `json:"done"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	FinishedAt   string `json:"finished_at"`
	SnoozedUntil string `json:"snoozed_until"`
	TodayUpdate  string `json:"todayUpdate"`
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

	currentDate, err := fromDateString(viewedDate)
	if err != nil {
		return ListTaskResponse{}, fmt.Errorf("invalid viewed date")
	}
	for _, task := range tasks {
		if truncate(task.CreatedAt).After(currentDate) {
			continue
		}
		if task.FinishedAt != nil {
			finishedAt := truncate(*task.FinishedAt)
			if finishedAt.Before(currentDate) {
				continue
			}
		}
		if task.SnoozedUntil != nil {
			snoozedUntil := truncate(*task.SnoozedUntil)
			if snoozedUntil.After(currentDate) {
				continue
			}
		}
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
				te.TodayUpdate = update.Description
			}
		}

		response.Tasks = append(response.Tasks, te)
	}

	return response, nil
}

type CreateTaskEntry struct {
	Title string `json:"title"`
	Date  string `json:"date"`
}

func (s *Service) CreateTask(newTaskEntry CreateTaskEntry) error {
	dateString, err := fromDateString(newTaskEntry.Date)
	if err != nil {
		return err
	}
	newTask := Task{
		Title:     newTaskEntry.Title,
		CreatedAt: dateString,
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
		return nil, fmt.Errorf("could not find entry with id %d", entry.Id)
	}

	t.Title = entry.Title

	s.Db.Save(t)

	return t, nil
}

type UpdateTaskDescriptionEntry struct {
	Id          uint   `json:"id"`
	Date        string `json:"date"`
	Description string `json:"description"`
}

func (s *Service) UpdateDailyUpdate(entry UpdateTaskDescriptionEntry) error {
	t := &Task{}
	result := s.Db.Model(Task{}).Preload("Updates").First(t, entry.Id)
	if result.Error != nil {
		return fmt.Errorf("could not find entry with id %d", entry.Id)
	}

	for _, update := range t.Updates {
		updateDate := toDateString(&update.Date)
		if updateDate == entry.Date {
			update.Description = entry.Description
			s.Db.Save(&update)
			return nil
		}
	}

	// Update for given date not found, create a new update

	updateDate, err := fromDateString(entry.Date)
	if err != nil {
		return fmt.Errorf("could not parse date %s", entry.Date)
	}
	newUpdate := Update{
		TaskId:      t.ID,
		Date:        updateDate,
		Description: entry.Description,
		DoneToday:   false,
	}
	return s.Db.Save(&newUpdate).Error
}

func (s *Service) DeleteTask(id string) error {
	return s.Db.Delete(&Task{}, id).Error
}

type SetTaskDoneEntry struct {
	Id   uint   `json:"id"`
	Date string `json:"date"`
	Done bool   `json:"done"`
}

func (s *Service) SetTaskDone(entry SetTaskDoneEntry) error {
	t := &Task{}
	result := s.Db.First(t, entry.Id)
	if result.Error != nil {
		return fmt.Errorf("could not find entry with id %d", entry.Id)
	}

	t.Done = entry.Done
	if t.Done == true {
		doneDate, err := fromDateString(entry.Date)
		if err != nil {
			return fmt.Errorf("could not parse date %s", entry.Date)
		}
		t.FinishedAt = &doneDate
	} else {
		t.FinishedAt = nil
	}

	return s.Db.Save(t).Error
}

type SnoozeTaskEntry struct {
	Id   uint   `json:"id"`
	Date string `json:"date"`
}

func (s *Service) SnoozeTask(entry SnoozeTaskEntry) error {
	t := &Task{}
	result := s.Db.First(t, entry.Id)
	if result.Error != nil {
		return fmt.Errorf("could not find entry with id %d", entry.Id)
	}

	snoozedDate, err := fromDateString(entry.Date)
	if err != nil {
		return fmt.Errorf("could not parse date %s", entry.Date)
	}
	t.SnoozedUntil = &snoozedDate

	return s.Db.Save(t).Error
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

func truncate(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
}
