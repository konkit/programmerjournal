package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"programmerjournal-backend/service"
)

func NewRouter(dbPath string) *mux.Router {
	h := New(dbPath)

	r := mux.NewRouter()
	r.HandleFunc("/api/tasks/list", h.ListTasksForDay)
	r.HandleFunc("/api/tasks/create", h.CreateTask)
	r.HandleFunc("/api/tasks/title/update", h.SetTaskTitle)
	r.HandleFunc("/api/tasks/description/update", h.SetTaskDailyUpdate)
	r.HandleFunc("/api/tasks/snooze", h.SnoozeTask)
	r.HandleFunc("/api/tasks/delete/{id}", h.DeleteTask)
	r.HandleFunc("/api/tasks/setDone", h.SetTaskDone)

	return r
}

type Handlers struct {
	s service.Service
}

func New(dbPath string) Handlers {
	s := service.New(dbPath)

	return Handlers{
		s: s,
	}
}

func (h *Handlers) ListTasksForDay(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		logAndWriteError("empty date query param", errors.New("empty date query param"), w)
		return
	}

	response, err := h.s.ListTasks(date)
	if err != nil {
		logAndWriteError("error listing tasks", err, w)
	}

	w.Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}

func (h *Handlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var createEntry service.CreateTaskEntry
	err := decoder.Decode(&createEntry)
	if err != nil {
		logAndWriteError("error decoding CreateTaskEntry", err, w)
		return
	}

	err = h.s.CreateTask(createEntry)
	if err != nil {
		logAndWriteError("error creating task", err, w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	//resp := SummaryEntry{Url: r.Url, Sum: votesCount}
	//_ = json.NewEncoder(w).Encode(resp)
}

func (h *Handlers) SetTaskTitle(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var createEntry service.UpdateTitleEntry
	err := decoder.Decode(&createEntry)
	if err != nil {
		logAndWriteError("error decoding UpdateTitleEntry", err, w)
		return
	}

	task, err := h.s.UpdateTaskTitle(createEntry)
	if err != nil {
		logAndWriteError("error creating task", err, w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(task)
}

func (h *Handlers) SetTaskDailyUpdate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var entry service.UpdateTaskDescriptionEntry
	err := decoder.Decode(&entry)
	if err != nil {
		logAndWriteError("error decoding UpdateTaskDescriptionEntry", err, w)
		return
	}

	err = h.s.UpdateDailyUpdate(entry)
	if err != nil {
		logAndWriteError("error updating task", err, w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) SetTaskDone(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var entry service.SetTaskDoneEntry
	err := decoder.Decode(&entry)
	if err != nil {
		logAndWriteError("error decoding SetTaskDoneEntry", err, w)
		return
	}

	err = h.s.SetTaskDone(entry)
	if err != nil {
		logAndWriteError("error updating task", err, w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	err := h.s.DeleteTask(vars["id"])
	if err != nil {
		logAndWriteError("error creating task", err, w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) SnoozeTask(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var entry service.SnoozeTaskEntry
	err := decoder.Decode(&entry)
	if err != nil {
		logAndWriteError("error decoding UpdateTitleEntry", err, w)
		return
	}

	err = h.s.SnoozeTask(entry)
	if err != nil {
		logAndWriteError("error creating task", err, w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

//func (h *Handlers) Root(w http.ResponseWriter, r *http.Request) {
//	res := fmt.Sprintf("%v", time.Now())
//	_, err := io.WriteString(w, res)
//	if err != nil {
//		log.Printf("Error constructing response: %v", err)
//		w.WriteHeader(http.StatusInternalServerError)
//		return
//	}
//}
//
//func (h *Handlers) GetVotesForUrl(w http.ResponseWriter, r *http.Request) {
//	url := r.URL.Query().Get("url")
//	if url == "" {
//		logAndWriteError("empty url query param", errors.New("empty url query param"), w)
//		return
//	}
//	sum, err := h.fetchVotesForUrl(url)
//	if err != nil {
//		logAndWriteError("error fetching votes: %v", err, w)
//		return
//	}
//	resp := SummaryEntry{Url: url, Sum: sum}
//	w.Header().Add("Content-Type", "application/json")
//	_ = json.NewEncoder(w).Encode(resp)
//}
//
//func (h *Handlers) VoteUp(w http.ResponseWriter, req *http.Request) {
//	decoder := json.NewDecoder(req.Body)
//
//	var r VoteRequest
//	err := decoder.Decode(&r)
//	if err != nil {
//		logAndWriteError("error decoding VoteRequest", err, w)
//		return
//	}
//
//	err = h.vote(r, 1)
//	if err != nil {
//		logAndWriteError("error handling vote up request", err, w)
//		return
//	}
//
//	votesCount, err := h.fetchVotesForUrl(r.Url)
//	if err != nil {
//		logAndWriteError("error fetching current votes count", err, w)
//		return
//	}
//
//	w.Header().Add("Content-Type", "application/json")
//	w.WriteHeader(http.StatusCreated)
//	resp := SummaryEntry{Url: r.Url, Sum: votesCount}
//	_ = json.NewEncoder(w).Encode(resp)
//}
//
//func (h *Handlers) vote(r VoteRequest, vote int) error {
//	stmt, err := h.db.Prepare("INSERT INTO votes(added, vote, url, comment) VALUES (?,?,?,?)")
//	if err != nil {
//		return fmt.Errorf("error inserting data into database: %v", err)
//	}
//
//	_, err = stmt.Exec(time.Now(), vote, r.Url, r.Comment)
//	if err != nil {
//		return fmt.Errorf("error inserting data into database: %v", err)
//	}
//
//	return nil
//}

func logAndWriteError(msg string, err error, w http.ResponseWriter) {
	log.Printf("%s: %v", msg, err)
	w.WriteHeader(http.StatusInternalServerError)
	_, e := io.WriteString(w, fmt.Sprintf("%s: %v", msg, err))
	if e != nil {
		log.Printf("error writing log: %v", e)
	}
	return
}
