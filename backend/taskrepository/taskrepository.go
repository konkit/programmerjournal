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
	err := r.db.Model(task.Task{}).
		Order("rank").
		Where("created_date = ?", date).
		Find(&tasksFromDB).
		Error

	return tasksFromDB, err
}

func (r *DBRepository) Create(newTask task.Task) error {
	var count int64
	r.db.Model(task.Task{}).Where("created_date = ?", newTask.CreatedDate).Count(&count)

	newTask.TaskID = uuid.NewString()
	newTask.Status = task.StatusCreated
	newTask.Update = ""
	newTask.Rank = int(count)

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
	t := &task.Task{ID: uint(taskID)}
	r.db.First(t)

	t.Title = title
	r.db.Save(&t)

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
	t := &task.Task{ID: uint(taskID)}
	r.db.First(t)

	t.Update = update
	r.db.Save(&t)

	return nil
}

func (r *DBRepository) GetTask(id uint64) (task.Task, error) {
	t := task.Task{ID: uint(id)}
	err := r.db.First(&t).Error
	return t, err
}

func (r *DBRepository) GetTasksByTaskID(taskID string) ([]task.Task, error) {
	var tasksFromDB []task.Task
	err := r.db.Model(task.Task{}).
		Where("task_id = ?", taskID).
		Find(&tasksFromDB).
		Error

	return tasksFromDB, err
}

func (r *DBRepository) ChangeRank(id uint64, newIndex int) error {
	t, err := r.GetTask(id)
	if err != nil {
		return err
	}

	// Find the reaction by ID
	oldIndex := t.Rank

	if oldIndex < newIndex {
		err = r.db.Model(&task.Task{}).
			Where("`rank` > ? AND `rank` <= ?", oldIndex, newIndex).
			Update("rank", gorm.Expr("`rank` - 1")).
			Error
	} else if oldIndex > newIndex {
		err = r.db.Model(&task.Task{}).
			Where("`rank` < ? AND `rank` >= ?", oldIndex, newIndex).
			Update("rank", gorm.Expr("`rank` + 1")).
			Error
	}

	// Update the order of the moved reaction
	t.Rank = newIndex
	err = r.db.Save(&t).Error
	if err != nil {
		return err
	}

	return err
}
