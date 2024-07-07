package main

import (
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

	r := handlers.NewRouter(dbPath)

	log.Printf("Starting server on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
