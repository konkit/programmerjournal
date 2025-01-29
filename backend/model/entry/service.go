package entry

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"programmerjournal-backend/model/date"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

func (r *Service) ListDayEntries(date date.DayDate) ([]Entry, error) {
	return r.listEntries(date.Value)
}

func (r *Service) ListMonthEntries(date date.MonthDate) ([]Entry, error) {
	return r.listEntries(date.Value)
}

func (r *Service) listEntries(date date.DateString) ([]Entry, error) {
	var entriesFromDB []Entry
	err := r.db.Model(Entry{}).
		Order("rank").
		Where("created_date = ?", date).
		Find(&entriesFromDB).
		Error

	return entriesFromDB, err
}

func (r *Service) WeeklySummary(firstDayOfWeek date.DayDate) (WeeklySummary, error) {
	isDateMonday := checkIfDateIsMonday(firstDayOfWeek)
	if !isDateMonday {
		return WeeklySummary{}, fmt.Errorf("the selected date is not the first day of the week")
	}

	monDate := firstDayOfWeek
	tueDate := monDate.PlusDays(1)
	wedDate := tueDate.PlusDays(1)
	thuDate := wedDate.PlusDays(1)
	friDate := thuDate.PlusDays(1)
	satDate := friDate.PlusDays(1)
	sunDate := satDate.PlusDays(1)

	var tasksFromDB []Entry
	err := r.db.Model(Entry{}).
		Where("created_date IN (?, ?, ?, ?, ?, ?, ?)", monDate.Value, tueDate.Value, wedDate.Value, thuDate.Value, friDate.Value, satDate.Value, sunDate.Value).
		Where("status LIKE ?", "Task%").
		Find(&tasksFromDB).
		Error

	if err != nil {
		return WeeklySummary{}, fmt.Errorf("failed to get weekly summary: %v", err)
	}

	taskMap := groupByTaskID(tasksFromDB)

	summaryArr := []TaskSummary{}
	for _, taskArr := range taskMap {
		t, err := findLastDayTask(taskArr)
		if err != nil {
			return WeeklySummary{}, err
		}
		updates := getUpdates(taskArr)

		summary := TaskSummary{
			TaskEntry: t,
			Updates:   updates,
		}
		summaryArr = append(summaryArr, summary)
	}

	notesFromDB := []Entry{}
	err = r.db.Model(Entry{}).
		Where("created_date IN (?, ?, ?, ?, ?, ?, ?)", monDate.Value, tueDate.Value, wedDate.Value, thuDate.Value, friDate.Value, satDate.Value, sunDate.Value).
		Where("status LIKE ?", "Note%").
		Find(&notesFromDB).
		Error

	if err != nil {
		return WeeklySummary{}, fmt.Errorf("failed to get weekly summary: %v", err)
	}

	ws := WeeklySummary{
		TaskSummaries: summaryArr,
		Notes:         notesFromDB,
	}

	return ws, err
}

func findLastDayTask(arr []Entry) (Entry, error) {
	lastTask := Entry{
		CreatedDate: "1000-01-01",
	}
	for _, t := range arr {
		createdDate, err := date.ParseDayDate(t.CreatedDate)
		if err != nil {
			return Entry{}, fmt.Errorf("entries dates not in day format: %s", t.CreatedDate)
		}
		lastTaskDate, err := date.ParseDayDate(lastTask.CreatedDate)
		if err != nil {
			return Entry{}, fmt.Errorf("entries dates not in day format: %s", lastTask.CreatedDate)
		}
		isAfter, err := createdDate.IsAfter(lastTaskDate)
		if err != nil {
			return Entry{}, err
		}

		if isAfter {
			lastTask = t
		}
	}

	return lastTask, nil
}

func getUpdates(arr []Entry) []TaskUpdate {
	var res []TaskUpdate
	for _, t := range arr {
		tu := TaskUpdate{
			Date:   t.CreatedDate,
			Update: t.TaskUpdate,
			Status: t.Status,
		}
		res = append(res, tu)
	}
	return res
}

func checkIfDateIsMonday(d date.DayDate) bool {
	// TODO: Check if monday
	return true
}

// Function to group a slice of structs by a specific property
func groupByTaskID(tasks []Entry) map[string][]Entry {
	grouped := make(map[string][]Entry)
	for _, tt := range tasks {
		grouped[tt.TaskID] = append(grouped[tt.TaskID], tt)
	}
	return grouped
}

func (r *Service) CreateTask(newTask Entry) error {
	nextRank := r.fetchNextRank(newTask.CreatedDate)

	newTask.TaskID = uuid.NewString()
	newTask.Status = StatusTaskCreated
	newTask.TaskUpdate = ""
	newTask.Rank = nextRank

	return r.db.Create(&newTask).Error
}

func (r *Service) CreateNote(newTask Entry) error {
	count := r.fetchNextRank(newTask.CreatedDate)

	newTask.TaskID = uuid.NewString()
	newTask.Status = StatusNote
	newTask.TaskUpdate = ""
	newTask.Rank = count

	return r.db.Create(&newTask).Error
}

func (r *Service) Update(entryID uint64, updatedTask Entry) error {
	updatedTask.ID = uint(entryID)
	return r.db.Save(updatedTask).Error
}

func (r *Service) Delete(entryID uint64) error {
	// TODO: Handle Rank change if the task is deleted in the middle
	return r.db.Delete(&Entry{}, entryID).Error
}

func (r *Service) SnoozeTask(entryID uint64, date date.DateString) error {
	snoozedTask, err := r.getEntryByID(entryID)
	if err != nil {
		return err
	}

	snoozeAfterTaskDate, err := date.IsAfter(snoozedTask.CreatedDate)
	if err != nil {
		return err
	}

	if !snoozeAfterTaskDate {
		return fmt.Errorf("snooze date must be in the future")
	}

	if snoozedTask.Status != StatusTaskCreated {
		return fmt.Errorf("invalid entry status, can only snooze created tasks")
	}

	snoozedTask.Status = StatusTaskSnoozed
	snoozedTask.TaskSnoozedUntil = date
	r.db.Save(&snoozedTask)

	nextRank := r.fetchNextRank(date)

	newTask := Clone(snoozedTask)
	newTask.Status = StatusTaskCreated
	newTask.CreatedDate = date
	newTask.Rank = nextRank
	r.db.Save(&newTask)

	return nil
}

func (r *Service) SetTitle(entryID uint64, title string) error {
	t, err := r.getEntryByID(entryID)
	if err != nil {
		return err
	}

	t.Title = title
	r.db.Save(&t)

	return nil
}

func (r *Service) SetTaskDone(entryID uint64, done bool) error {
	t, err := r.getEntryByID(entryID)
	if err != nil {
		return err
	}

	if done == true {
		t.Status = StatusTaskDone
	} else {
		t.Status = StatusTaskCreated
	}
	r.db.Save(&t)

	return nil
}

func (r *Service) MigrateToDaily(entryID uint64, date date.DayDate) error {
	t, err := r.getEntryByID(entryID)
	if err != nil {
		return err
	}

	t.Status = StatusTaskMigrated
	t.TaskSnoozedUntil = date.Value
	r.db.Save(&t)

	nextRank := r.fetchNextRank(date.Value)

	newTask := Clone(t)
	newTask.Status = StatusTaskCreated
	newTask.CreatedDate = date.Value
	newTask.Rank = nextRank
	r.db.Save(&newTask)

	return nil
}

func (r *Service) MigrateToMonthly(entryID uint64, date date.MonthDate) error {
	t, err := r.getEntryByID(entryID)
	if err != nil {
		return err
	}

	t.Status = StatusTaskMigrated
	t.TaskSnoozedUntil = date.Value
	r.db.Save(&t)

	nextRank := r.fetchNextRank(date.Value)

	newTask := Clone(t)
	newTask.Status = StatusTaskCreated
	newTask.CreatedDate = date.Value
	newTask.Rank = nextRank
	r.db.Save(&newTask)

	return nil
}

func (r *Service) ImportPastTasksFromDay(today date.DayDate) error {
	for i := 1; i < 30; i++ {
		current := today.MinusDays(i)

		tasks, err := r.ListDayEntries(current)
		if err != nil {
			return err
		}

		for _, t := range tasks {
			if t.Status == StatusTaskCreated {
				nextRank := r.fetchNextRank(today.Value)

				newTask := Clone(t)
				newTask.CreatedDate = today.Value
				newTask.Rank = nextRank
				newTask.TaskUpdate = ""
				r.db.Save(&newTask)

				t.Status = StatusTaskSnoozed
				r.db.Save(&t)
			}
		}
	}

	return nil
}

func (r *Service) ImportPastTasksFromMonth(today date.MonthDate) error {
	for i := 1; i < 12; i++ {
		current := today.MinusMonth(i)

		tasks, err := r.ListMonthEntries(current)
		if err != nil {
			return err
		}

		for _, t := range tasks {
			if t.Status == StatusTaskCreated {
				nextRank := r.fetchNextRank(today.Value)

				newTask := Clone(t)
				newTask.CreatedDate = today.Value
				newTask.Rank = nextRank
				newTask.TaskUpdate = ""
				r.db.Save(&newTask)

				t.Status = StatusTaskSnoozed
				r.db.Save(&t)
			}
		}
	}

	return nil
}

func (r *Service) SetTaskUpdate(entryID uint64, update string) error {
	t, err := r.getEntryByID(entryID)
	if err != nil {
		return err
	}

	t.TaskUpdate = update
	r.db.Save(&t)

	return nil
}

func (r *Service) SetDescription(entryID uint64, description string) error {
	t, err := r.getEntryByID(entryID)
	if err != nil {
		return err
	}

	t.Description = description
	r.db.Save(&t)

	return nil
}

func (r *Service) ChangeRank(entryID uint64, newIndex int) error {
	e, err := r.getEntryByID(entryID)
	if err != nil {
		return err
	}
	oldIndex := e.Rank

	var entriesFromDB []Entry
	err = r.db.Model(Entry{}).
		Order("rank").
		Where("created_date = ?", e.CreatedDate).
		Find(&entriesFromDB).
		Error

	var entriesWithoutMovedElement []Entry
	elementFound := false
	for _, entry := range entriesFromDB {
		if entry.Rank != oldIndex {
			entriesWithoutMovedElement = append(entriesWithoutMovedElement, entry)
		} else if elementFound == true {
			entriesWithoutMovedElement = append(entriesWithoutMovedElement, entry)
		} else {
			elementFound = true
		}
	}

	//if len(entriesWithoutMovedElement) == 0 {
	//	fmt.Println("entriesWithoutMovedElement is empty")
	//	return nil
	//}

	elementAdded := false
	var currentRank int
	if len(entriesWithoutMovedElement) > 0 {
		currentRank = min(entriesWithoutMovedElement[0].Rank, newIndex)
		entriesIter := 0

		for entriesIter < len(entriesWithoutMovedElement) {
			if currentRank == newIndex {
				err = r.saveWithNewRank(e, currentRank)
				if err != nil {
					return err
				}
				elementAdded = true
			} else {
				entryToAdd := entriesWithoutMovedElement[entriesIter]
				if currentRank >= 0 || entryToAdd.Rank < 0 {
					entriesIter++
					err = r.saveWithNewRank(entryToAdd, currentRank)
					if err != nil {
						return err
					}
				}
			}

			currentRank++
		}
	} else {
		currentRank = min(newIndex, 0)
	}

	if !elementAdded {
		if currentRank < 0 && newIndex >= 0 {
			currentRank = 0
		}
		err = r.saveWithNewRank(e, currentRank)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *Service) saveWithNewRank(e Entry, currentRank int) error {
	// Do not save if the rank is already set to currentRank
	if e.Rank != currentRank {
		e.Rank = currentRank
		return r.db.Save(e).Error
	}
	return nil
}

func (r *Service) GetTaskSummary(id uint64) (*TaskSummary, error) {
	e, err := r.getEntryByID(id)
	if err != nil {
		return nil, err
	}

	var tasksFromDB []Entry
	err = r.db.Model(Entry{}).
		Where("task_id = ?", e.TaskID).
		Find(&tasksFromDB).
		Error

	if err != nil {
		return nil, err
	}

	var updates []TaskUpdate
	for _, tt := range tasksFromDB {
		update := TaskUpdate{
			Date:   tt.CreatedDate,
			Update: tt.TaskUpdate,
			Status: tt.Status,
		}
		updates = append(updates, update)
	}

	ts := &TaskSummary{
		TaskEntry: e,
		Updates:   updates,
	}

	return ts, nil
}

func (r *Service) getEntryByID(entryID uint64) (Entry, error) {
	t := Entry{ID: uint(entryID)}
	if err := r.db.First(&t).Error; err != nil {
		return t, err
	}
	return t, nil
}

func (r *Service) fetchNextRank(date date.DateString) int {
	var count int64
	r.db.Model(Entry{}).
		Where("created_date = ?", date).
		Where("rank >= 0").
		Count(&count)
	return int(count)
}
