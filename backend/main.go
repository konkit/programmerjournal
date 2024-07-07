package main

import (
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"programmerjournal-backend/handlers"
)

func main() {
	var dbPath = "./foo.db"
	var port = "8080"
	if len(os.Args) >= 2 {
		dbPath = os.Args[1]
	}
	if len(os.Args) >= 3 {
		port = os.Args[2]
	}

	h := handlers.New(dbPath)

	r := mux.NewRouter()
	r.HandleFunc("/api/tasks/list", h.ListTasksForDay)
	r.HandleFunc("/api/tasks/create", h.CreateTask)
	r.HandleFunc("/api/tasks/title/update", h.SetTaskTitle)

	log.Printf("Starting server on port %s", port)

	log.Fatal(http.ListenAndServe(":"+port, r))
}
