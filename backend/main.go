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

	db := handlers.InitDB(dbPath)
	r := handlers.NewRouter(db)

	log.Printf("Starting server on port %s", port)
	listenStr := ":" + port
	log.Fatal(http.ListenAndServe(listenStr, r))
}
