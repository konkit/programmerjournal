package taskrepository

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"programmerjournal-backend/date"
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

func (r *DBRepository) ListTasks(date date.Date) ([]task.Task, error) {
	var tasksFromDB []task.Task
	err := r.db.Model(task.Task{}).
		Order("rank").
		Where("created_date = ?", date).
		Find(&tasksFromDB).
		Error

	return tasksFromDB, err
}

func (r *DBRepository) WeeklySummary(firstDayOfWeek date.Date) ([]task.TaskWeeklySummary, error) {
	isDateMonday := checkIfDateIsMonday(firstDayOfWeek)
	if !isDateMonday {
		return nil, fmt.Errorf("the selected date is not the first day of the week")
	}

	monDate := firstDayOfWeek
	tueDate := date.PlusDays(monDate, 1)
	wedDate := date.PlusDays(tueDate, 1)
	thuDate := date.PlusDays(wedDate, 1)
	friDate := date.PlusDays(thuDate, 1)
	satDate := date.PlusDays(friDate, 1)
	sunDate := date.PlusDays(satDate, 1)

	var tasksFromDB []task.Task
	err := r.db.Model(task.Task{}).
		Where("created_date IN (?, ?, ?, ?, ?, ?, ?)", monDate, tueDate, wedDate, thuDate, friDate, satDate, sunDate).
		Find(&tasksFromDB).
		Error

	if err != nil {
		return nil, fmt.Errorf("failed to get weekly summary: %v", err)
	}

	taskMap := groupByTaskID(tasksFromDB)

	var summaryArr []task.TaskWeeklySummary

	for _, taskArr := range taskMap {
		t, err := findLastTask(taskArr)
		if err != nil {
			return nil, err
		}
		updates := getUpdates(taskArr)

		summary := task.TaskWeeklySummary{
			Task:    t,
			Updates: updates,
		}
		summaryArr = append(summaryArr, summary)
	}

	return summaryArr, err
}

func findLastTask(arr []task.Task) (task.Task, error) {
	lastTask := task.Task{
		CreatedDate: "1000-01-01",
	}
	for _, t := range arr {
		isAfter, err := t.CreatedDate.IsAfter(lastTask.CreatedDate)
		if err != nil {
			return task.Task{}, err
		}

		if isAfter {
			lastTask = t
		}
	}

	return lastTask, nil
}

func getUpdates(arr []task.Task) []task.TaskUpdate {
	var res []task.TaskUpdate
	for _, t := range arr {
		tu := task.TaskUpdate{
			Date:   t.CreatedDate,
			Update: t.Update,
			Status: t.Status,
		}
		res = append(res, tu)
	}
	return res
}

func checkIfDateIsMonday(d date.Date) bool {
	// TODO: Check if monday
	return true
}

// Function to group a slice of structs by a specific property
func groupByTaskID(tasks []task.Task) map[string][]task.Task {
	grouped := make(map[string][]task.Task)
	for _, tt := range tasks {
		grouped[tt.TaskID] = append(grouped[tt.TaskID], tt)
	}
	return grouped
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

func (r *DBRepository) CreateNote(newTask task.Task) error {
	var count int64
	r.db.Model(task.Task{}).Where("created_date = ?", newTask.CreatedDate).Count(&count)

	newTask.TaskID = uuid.NewString()
	newTask.Status = task.StatusNote
	newTask.Update = ""
	newTask.Rank = int(count)

	return r.db.Create(&newTask).Error
}

func (r *DBRepository) Update(taskID uint64, updatedTask task.Task) error {
	updatedTask.ID = uint(taskID)
	return r.db.Save(updatedTask).Error
}

func (r *DBRepository) Delete(taskID uint64) error {
	return r.db.Delete(&task.Task{}, taskID).Error
}

func (r *DBRepository) Snooze(taskID uint64, date date.Date) error {
	snoozedTask, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	snoozedTask.Status = task.StatusSnoozed
	snoozedTask.SnoozedUntil = date
	r.db.Save(&snoozedTask)

	var count int64
	r.db.Model(task.Task{}).Where("created_date = ?", date).Count(&count)

	newTask := task.Clone(snoozedTask)
	newTask.Status = task.StatusCreated
	newTask.CreatedDate = date
	newTask.Rank = int(count)
	r.db.Save(&newTask)

	return nil
}

func (r *DBRepository) SetTaskTitle(taskID uint64, title string) error {
	t, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	t.Title = title
	r.db.Save(&t)

	return nil
}

func (r *DBRepository) SetTaskDone(taskID uint64, done bool) error {
	t, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	if done == true {
		t.Status = task.StatusDone
	} else {
		t.Status = task.StatusCreated
	}
	r.db.Save(&t)

	return nil
}

func (r *DBRepository) MigrateToMonthly(taskID uint64, date date.Date) error {
	t, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	t.Status = task.StatusMigrated
	t.SnoozedUntil = date
	r.db.Save(&t)

	var count int64
	r.db.Model(task.Task{}).Where("created_date = ?", date).Count(&count)

	newTask := task.Clone(t)
	newTask.Status = task.StatusCreated
	newTask.CreatedDate = date
	newTask.Rank = int(count)
	r.db.Save(&newTask)

	return nil
}

func (r *DBRepository) ImportPastTasks(today date.Date) error {
	for i := 1; i < 30; i++ {
		current := date.MinusDays(today, i)

		tasks, err := r.ListTasks(current)
		if err != nil {
			return err
		}

		for _, t := range tasks {
			if t.Status == task.StatusCreated {
				var count int64
				r.db.Model(task.Task{}).Where("created_date = ?", today).Count(&count)

				newTask := task.Clone(t)
				newTask.CreatedDate = today
				newTask.Rank = int(count)
				r.db.Save(&newTask)

				t.Status = task.StatusSnoozed
				r.db.Save(&t)
			}
		}
	}

	return nil
}

func (r *DBRepository) SetTaskUpdate(taskID uint64, update string) error {
	t, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	t.Update = update
	r.db.Save(&t)

	return nil
}

func (r *DBRepository) SetTaskDescription(taskID uint64, description string) error {
	t, err := r.getTaskByID(taskID)
	if err != nil {
		return err
	}

	t.Description = description
	r.db.Save(&t)

	return nil
}

func (r *DBRepository) ChangeRank(id uint64, newIndex int) error {
	t, err := r.getTaskByID(id)
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

func (r *DBRepository) GetTaskSummary(id uint64) (*task.TaskSummary, error) {
	t, err := r.getTaskByID(id)
	if err != nil {
		return nil, err
	}

	var tasksFromDB []task.Task
	err = r.db.Model(task.Task{}).
		Where("task_id = ?", t.TaskID).
		Find(&tasksFromDB).
		Error

	if err != nil {
		return nil, err
	}

	var updates []task.TaskUpdate
	for _, tt := range tasksFromDB {
		update := task.TaskUpdate{
			Date:   tt.CreatedDate,
			Update: tt.Update,
			Status: tt.Status,
		}
		updates = append(updates, update)
	}

	ts := &task.TaskSummary{
		Task:    t,
		Updates: updates,
	}

	return ts, nil
}

func (r *DBRepository) getTaskByID(taskID uint64) (task.Task, error) {
	t := task.Task{ID: uint(taskID)}
	if err := r.db.First(&t).Error; err != nil {
		return t, err
	}
	return t, nil
}
