package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strconv"
)

type TaskStatus string

const (
	TaskCreated   TaskStatus = "Created"
	TaskDone                 = "Done"
	TaskSnoozed              = "Snoozed"
	TaskCancelled            = "Cancelled"
)

type Task struct {
	ID          uint       `json:"id" sql:"AUTO_INCREMENT" gorm:"primarykey"`
	TaskID      string     `json:"taskID"`
	Title       string     `json:"title"`
	Status      TaskStatus `json:"status"`
	CreatedDate string     `json:"createdDate"`
	Update      string     `json:"update"`
}

func InitDB(dbPath string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Task{})
	return db
}

func NewRouter(db *gorm.DB) *mux.Router {
	h := Handlers{db: db}

	r := mux.NewRouter()
	r.HandleFunc("/api/tasks/list/{date}", h.ListTasks)
	r.HandleFunc("/api/tasks/create", h.CreateTask)
	r.HandleFunc("/api/tasks/{id}/update", h.UpdateTask)
	r.HandleFunc("/api/tasks/{id}/delete", h.DeleteTask)
	r.HandleFunc("/api/tasks/{id}/snooze", h.SnoozeTask)
	r.HandleFunc("/api/tasks/{id}/setDone", h.SetTaskDone)
	r.HandleFunc("/api/tasks/{id}/setTitle", h.SetTaskTitle)
	r.HandleFunc("/api/tasks/{id}/setUpdate", h.SetTaskUpdate)
	r.HandleFunc("/api/tasks/summary/{taskID}", h.GetTask)
	r.HandleFunc("/api/tasks/moveToNextDay", h.MoveTasksToNextDay)

	return r
}

type Handlers struct {
	db *gorm.DB
}

func NewHandlers(db *gorm.DB) Handlers {
	return Handlers{db: db}
}

func (h *Handlers) ListTasks(w http.ResponseWriter, r *http.Request) {
	date, err := getDateFromParam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logAndWriteError("missing param 'date'", err, w)
		return
	}

	tasks, err := h.getTasksFromDate(date)
	if err != nil {
		logAndWriteError("error fetching tasks from database", err, w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tasks)
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		logAndWriteError("error decoding CreateTaskEntry", err, w)
		return
	}

	newTask.TaskID = uuid.NewString()
	newTask.Status = TaskCreated
	newTask.Update = ""

	err = h.db.Create(&newTask).Error
	if err != nil {
		logAndWriteError("Error saving task to the database", err, w)
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := getIDFromParam(r)
	if err != nil {
		logAndWriteError("error decoding ID", err, w)
		return
	}

	var updatedTask Task
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		logAndWriteError("error decoding UpdateTitleEntry", err, w)
		return
	}
	updatedTask.ID = uint(taskID)

	h.db.Save(updatedTask)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := getIDFromParam(r)
	if err != nil {
		logAndWriteError("error decoding ID", err, w)
		return
	}

	err = h.db.Delete(&Task{}, taskID).Error
	if err != nil {
		logAndWriteError("error creating task", err, w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

type SnoozeTaskEntry struct {
	Date string `json:"date"`
}

// SnoozeTask creates a related entry in the future time
func (h *Handlers) SnoozeTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := getIDFromParam(r)
	if err != nil {
		logAndWriteError("error decoding ID", err, w)
		return
	}

	var entry SnoozeTaskEntry
	err = json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		logAndWriteError("error decoding UpdateTitleEntry", err, w)
		return
	}

	snoozedTask := Task{ID: uint(taskID)}
	h.db.First(&snoozedTask)

	snoozedTask.Status = TaskSnoozed
	h.db.Save(&snoozedTask)

	newTask := cloneTask(snoozedTask)
	newTask.Status = TaskCreated
	newTask.CreatedDate = entry.Date
	h.db.Save(&newTask)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

type SetTaskDoneEntry struct {
	Done bool `json:"done"`
}

// SetTaskDone marks a task as done
func (h *Handlers) SetTaskDone(w http.ResponseWriter, r *http.Request) {
	taskID, err := getIDFromParam(r)
	if err != nil {
		logAndWriteError("error decoding ID", err, w)
		return
	}

	var entry SetTaskDoneEntry
	err = json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		logAndWriteError("error decoding SetTaskDoneEntry", err, w)
		return
	}

	task := &Task{ID: uint(taskID)}
	h.db.First(task)

	if entry.Done == true {
		task.Status = TaskDone
	} else {
		task.Status = TaskCreated
	}
	h.db.Save(&task)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

type SetTaskTitleEntry struct {
	Title string `json:"title"`
}

// SetTaskTitle changes title for the entry of a given day's entry
func (h *Handlers) SetTaskTitle(w http.ResponseWriter, r *http.Request) {
	taskID, err := getIDFromParam(r)
	if err != nil {
		logAndWriteError("error decoding ID", err, w)
		return
	}

	var entry SetTaskTitleEntry
	err = json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		logAndWriteError("error decoding SetTaskDoneEntry", err, w)
		return
	}

	task := &Task{ID: uint(taskID)}
	h.db.First(task)

	task.Title = entry.Title
	h.db.Save(&task)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

type SetTaskUpdateEntry struct {
	Update string `json:"update"`
}

// SetTaskUpdate changes update for the entry of a given day's entry
func (h *Handlers) SetTaskUpdate(w http.ResponseWriter, r *http.Request) {
	taskID, err := getIDFromParam(r)
	if err != nil {
		logAndWriteError("error decoding ID", err, w)
		return
	}

	var entry SetTaskUpdateEntry
	err = json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		logAndWriteError("error decoding SetTaskDoneEntry", err, w)
		return
	}

	task := &Task{ID: uint(taskID)}
	h.db.First(task)

	task.Update = entry.Update
	h.db.Save(&task)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

type MoveTasksToNextDayEntry struct {
	Current string `json:"current"`
	Next    string `json:"next"`
}

// MoveTasksToNextDay moves tasks to the given day, creating new Task entries, but preserving their TaskID number,
// so that they are later discoverable.
func (h *Handlers) MoveTasksToNextDay(w http.ResponseWriter, r *http.Request) {
	var entry MoveTasksToNextDayEntry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		logAndWriteError("error decoding MoveTasksToNextDayEntry", err, w)
		return
	}

	tasks, err := h.getTasksFromDate(entry.Current)
	if err != nil {
		logAndWriteError("error fetching tasks from database", err, w)
	}

	for _, task := range tasks {
		if task.Status == TaskCreated {
			newTask := cloneTask(task)
			newTask.CreatedDate = entry.Next
			h.db.Save(&newTask)

			task.Status = TaskSnoozed
			h.db.Save(&task)
		}
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

// GetTask gets a task from a day, but also lists all updates from past days
func (h *Handlers) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := getTaskIDFromParam(r)
	if err != nil {
		logAndWriteError("error decoding ID", err, w)
		return
	}

	tasks, err := h.getTasksByTaskID(taskID)
	if err != nil {
		logAndWriteError("error fetching tasks from database", err, w)
	}

	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tasks)
}

func (h *Handlers) getTasksFromDate(date string) ([]Task, error) {
	var tasksFromDB []Task
	err := h.db.Model(Task{}).Find(&tasksFromDB).Error

	var filteredTasks []Task = []Task{}
	for _, task := range tasksFromDB {
		if task.CreatedDate != date {
			continue
		}

		filteredTasks = append(filteredTasks, task)
	}

	return filteredTasks, err
}

func (h *Handlers) getTasksByTaskID(taskID string) ([]Task, error) {
	var tasksFromDB []Task
	err := h.db.Model(Task{}).Find(&tasksFromDB).Error

	var filteredTasks []Task = []Task{}
	for _, task := range tasksFromDB {
		if task.TaskID != taskID {
			continue
		}

		filteredTasks = append(filteredTasks, task)
	}

	return filteredTasks, err
}

func logAndWriteError(msg string, err error, w http.ResponseWriter) {
	log.Printf("%s: %v", msg, err)
	w.WriteHeader(http.StatusInternalServerError)
	_, e := io.WriteString(w, fmt.Sprintf("%s: %v", msg, err))
	if e != nil {
		log.Printf("error writing log: %v", e)
	}
	return
}

func getIDFromParam(r *http.Request) (uint64, error) {
	vars := mux.Vars(r)
	taskIDStr := vars["id"]

	taskID, err := strconv.ParseUint(taskIDStr, 10, 0)
	return taskID, err
}

func getTaskIDFromParam(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	taskIDStr := vars["taskID"]
	if taskIDStr == "" {
		return "", fmt.Errorf("empty param 'taskID'")
	}
	return taskIDStr, nil
}

func getDateFromParam(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	date, ok := vars["date"]
	if !ok {
		return "", fmt.Errorf("missing param date")
	}

	if date == "" {
		return "", fmt.Errorf("missing param date")
	}

	return date, nil
}

func cloneTask(src Task) Task {
	newTask := Task{
		Title:       src.Title,
		TaskID:      src.TaskID,
		Status:      src.Status,
		Update:      src.Update,
		CreatedDate: src.CreatedDate,
	}
	return newTask
}
