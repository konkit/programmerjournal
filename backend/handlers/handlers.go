package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
	"strconv"
)

func NewRouter(dbRepo *taskrepository.DBRepository) *mux.Router {
	h := Handlers{taskRepo: dbRepo}

	r := mux.NewRouter()
	r.HandleFunc("/api/tasks/list/{date}", h.ListTasks)
	r.HandleFunc("/api/tasks/create", h.CreateTask)
	r.HandleFunc("/api/tasks/{id}/update", h.UpdateTask)
	r.HandleFunc("/api/tasks/{id}/summary", h.GetTaskSummary)
	r.HandleFunc("/api/tasks/{id}/delete", h.DeleteTask)
	r.HandleFunc("/api/tasks/{id}/snooze", h.SnoozeTask)
	r.HandleFunc("/api/tasks/{id}/setDone", h.SetTaskDone)
	r.HandleFunc("/api/tasks/{id}/setTitle", h.SetTaskTitle)
	r.HandleFunc("/api/tasks/{id}/setUpdate", h.SetTaskUpdate)
	r.HandleFunc("/api/tasks/{id}/setDescription", h.SetTaskDescription)
	//r.HandleFunc("/api/tasks/summary/{taskID}", h.GetTaskSummary)
	r.HandleFunc("/api/tasks/moveToNextDay", h.MoveTasksToNextDay)
	r.HandleFunc("/api/tasks/updateFromPastDays", h.UpdateFromPastDays)
	r.HandleFunc("/api/tasks/{id}/changeRank", h.ChangeRank)
	r.HandleFunc("/api/tasks/importPastTasks/{date}", h.ImportPastTasks)

	return r
}

type Handlers struct {
	taskRepo *taskrepository.DBRepository
}

func NewHandlers(dbRepo *taskrepository.DBRepository) Handlers {
	return Handlers{taskRepo: dbRepo}
}

func (h *Handlers) ListTasks(w http.ResponseWriter, r *http.Request) {
	date, err := getDateFromParam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logAndWriteError("missing param 'date'", err, w)
		return
	}

	tasks, err := h.taskRepo.GetTasksFromDate(date)
	if err != nil {
		logAndWriteError("error fetching tasks from database", err, w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(tasks)
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var newTask task.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		logAndWriteError("error decoding CreateTaskEntry", err, w)
		return
	}

	err = h.taskRepo.Create(newTask)
	if err != nil {
		logAndWriteError("Error saving task to the database", err, w)
	}

	writeResponseCreated(w)
}

func (h *Handlers) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := getIDFromParam(r)
	if err != nil {
		logAndWriteError("error decoding ID", err, w)
		return
	}

	var updatedTask task.Task
	err = json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		logAndWriteError("error decoding UpdateTitleEntry", err, w)
		return
	}

	h.taskRepo.Update(taskID, updatedTask)

	writeResponseOK(w)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := getIDFromParam(r)
	if err != nil {
		logAndWriteError("error decoding ID", err, w)
		return
	}

	err = h.taskRepo.Delete(taskID)
	if err != nil {
		logAndWriteError("error creating task", err, w)
		return
	}

	writeResponseOK(w)
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

	err = h.taskRepo.Snooze(taskID, entry.Date)
	if err != nil {
		logAndWriteError("error snoozing task", err, w)
	}

	writeResponseOK(w)
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

	err = h.taskRepo.SetTaskDone(taskID, entry.Done)
	if err != nil {
		logAndWriteError("error setting task as done", err, w)
		return
	}

	writeResponseOK(w)
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

	err = h.taskRepo.SetTaskTitle(taskID, entry.Title)
	if err != nil {
		logAndWriteError("error setting task as title", err, w)
		return
	}

	writeResponseOK(w)
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

	err = h.taskRepo.SetTaskUpdate(taskID, entry.Update)
	if err != nil {
		logAndWriteError("error setting update", err, w)
		return
	}

	writeResponseOK(w)
}

type SetTaskDescriptionEntry struct {
	Description string `json:"description"`
}

// SetTaskDescription changes description for the entry of a given task
func (h *Handlers) SetTaskDescription(w http.ResponseWriter, r *http.Request) {
	taskID, err := getIDFromParam(r)
	if err != nil {
		logAndWriteError("error decoding ID", err, w)
		return
	}

	var entry SetTaskDescriptionEntry
	err = json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		logAndWriteError("error decoding SetTaskDoneEntry", err, w)
		return
	}

	err = h.taskRepo.SetTaskDescription(taskID, entry.Description)
	if err != nil {
		logAndWriteError("error setting update", err, w)
		return
	}

	writeResponseOK(w)
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

	err = h.taskRepo.MoveTasksToNextDay(entry.Current, entry.Next)
	if err != nil {
		logAndWriteError("error moving tasks to next day", err, w)
		return
	}
	writeResponseOK(w)
}

// MoveTasksToNextDay moves tasks to the given day, creating new Task entries, but preserving their TaskID number,
// so that they are later discoverable.
func (h *Handlers) ImportPastTasks(w http.ResponseWriter, r *http.Request) {
	date, err := getDateFromParam(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logAndWriteError("missing param 'date'", err, w)
		return
	}

	err = h.taskRepo.ImportPastTasks(date)
	if err != nil {
		logAndWriteError("error moving tasks to next day", err, w)
		return
	}
	writeResponseOK(w)
}

// GetTaskSummary gets a task from a day, but also lists all updates from past days
func (h *Handlers) GetTaskSummary(w http.ResponseWriter, r *http.Request) {
	id, err := getIDFromParam(r)
	if err != nil {
		logAndWriteError("error decoding ID", err, w)
		return
	}

	ts, err := h.taskRepo.GetTaskSummary(id)
	if err != nil {
		logAndWriteError("error fetching tasks from database", err, w)
	}

	//taskID, err := getTaskIDFromParam(r)
	//if err != nil {
	//	logAndWriteError("error decoding ID", err, w)
	//	return
	//}

	//tasks, err := h.taskRepo.GetTasksByTaskID(taskID)
	//if err != nil {
	//	logAndWriteError("error fetching tasks from database", err, w)
	//}

	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(ts)
}

func (h *Handlers) UpdateFromPastDays(writer http.ResponseWriter, request *http.Request) {

}

type ChangeRankEntry struct {
	NewRank int `json:"newRank"`
}

func (h *Handlers) ChangeRank(w http.ResponseWriter, r *http.Request) {
	var entry ChangeRankEntry
	err := json.NewDecoder(r.Body).Decode(&entry)
	if err != nil {
		logAndWriteError("error decoding ChangeRankEntry", err, w)
		return
	}

	id, err := getIDFromParam(r)
	if err != nil {
		logAndWriteError("error decoding ID", err, w)
		return
	}

	err = h.taskRepo.ChangeRank(id, entry.NewRank)
	if err != nil {
		logAndWriteError("error changing rank", err, w)
		return
	}

	writeResponseOK(w)
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

func logAndWriteError(msg string, err error, w http.ResponseWriter) {
	log.Printf("%s: %v", msg, err)
	w.WriteHeader(http.StatusInternalServerError)
	_, e := io.WriteString(w, fmt.Sprintf("%s: %v", msg, err))
	if e != nil {
		log.Printf("error writing log: %v", e)
	}
	return
}

func writeResponseOK(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func writeResponseCreated(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}
