package taskrepository

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"programmerjournal-backend/task"
)

func InitDB(dbPath string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	err = db.AutoMigrate(&task.Task{})
	if err != nil {
		return nil, fmt.Errorf("failed to migrate database")
	}
	return db, nil
}

type DBRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) (*DBRepository, error) {
	return &DBRepository{db: db}, nil
}

func (r *DBRepository) GetTasksFromDate(date string) ([]task.Task, error) {
	var tasksFromDB []task.Task
	err := r.db.Model(task.Task{}).Find(&tasksFromDB).Error

	var filteredTasks []task.Task = []task.Task{}
	for _, task := range tasksFromDB {
		if task.CreatedDate != date {
			continue
		}

		filteredTasks = append(filteredTasks, task)
	}

	return filteredTasks, err
}

func (r *DBRepository) Create(newTask task.Task) error {
	newTask.TaskID = uuid.NewString()
	newTask.Status = task.StatusCreated
	newTask.Update = ""

	return r.db.Create(&newTask).Error
}

func (r *DBRepository) Update(taskID uint64, updatedTask task.Task) {
	updatedTask.ID = uint(taskID)
	r.db.Save(updatedTask)
}

func (r *DBRepository) Delete(taskID uint64) error {
	return r.db.Delete(&task.Task{}, taskID).Error
}

func (r *DBRepository) Snooze(taskID uint64, date string) error {
	snoozedTask := task.Task{ID: uint(taskID)}
	r.db.First(&snoozedTask)

	snoozedTask.Status = task.StatusSnoozed
	r.db.Save(&snoozedTask)

	newTask := task.Clone(snoozedTask)
	newTask.Status = task.StatusCreated
	newTask.CreatedDate = date
	r.db.Save(&newTask)

	return nil
}

func (r *DBRepository) SetTaskTitle(taskID uint64, title string) error {
	task := &task.Task{ID: uint(taskID)}
	r.db.First(task)

	task.Title = title
	r.db.Save(&task)

	return nil
}

func (r *DBRepository) SetTaskDone(taskID uint64, done bool) error {
	t := &task.Task{ID: uint(taskID)}
	r.db.First(t)

	if done == true {
		t.Status = task.StatusDone
	} else {
		t.Status = task.StatusCreated
	}
	r.db.Save(&t)

	return nil
}

func (r *DBRepository) MoveTasksToNextDay(current string, next string) error {
	tasks, err := r.GetTasksFromDate(current)
	if err != nil {
		return err
	}

	for _, t := range tasks {
		if t.Status == task.StatusCreated {
			newTask := task.Clone(t)
			newTask.CreatedDate = next
			r.db.Save(&newTask)

			t.Status = task.StatusSnoozed
			r.db.Save(&t)
		}
	}

	return nil
}

func (r *DBRepository) SetTaskUpdate(taskID uint64, update string) error {
	task := &task.Task{ID: uint(taskID)}
	r.db.First(task)

	task.Update = update
	r.db.Save(&task)

	return nil
}

func (r *DBRepository) GetTasksByTaskID(taskID string) ([]task.Task, error) {
	var tasksFromDB []task.Task
	err := r.db.Model(task.Task{}).Find(&tasksFromDB).Error

	var filteredTasks []task.Task = []task.Task{}
	for _, task := range tasksFromDB {
		if task.TaskID != taskID {
			continue
		}

		filteredTasks = append(filteredTasks, task)
	}

	return filteredTasks, err
}
