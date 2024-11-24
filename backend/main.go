package main

import (
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
	"os"
	"programmerjournal-backend/handlers"
	"programmerjournal-backend/taskrepository"
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

	db, err := taskrepository.InitDB(dbPath)
	if err != nil {
		panic(err)
	}
	dbRepo, err := taskrepository.NewRepository(db)
	if err != nil {
		panic(err)
	}
	r := handlers.NewRouter(dbRepo)

	log.Printf("Starting server on port %s", port)
	listenStr := ":" + port
	log.Fatal(http.ListenAndServe(listenStr, r))
}
