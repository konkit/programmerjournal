package main

import (
	"embed"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"io/fs"
	"log"
	"net/http"
	"programmerjournal-backend/model/database"
	"strings"

	entryhandlers "programmerjournal-backend/handlers/entry"
	notehandlers "programmerjournal-backend/handlers/note"
	recurringtaskhandlers "programmerjournal-backend/handlers/recurringtask"
	taskhandlers "programmerjournal-backend/handlers/task"
)

//go:embed static
var staticFiles embed.FS

// Options for the CLI. Pass `--port` or set the `SERVICE_PORT` env var.
type Options struct {
	DBPath string `help:"Path to the database file" default:"./foo.db"`
	Port   int    `help:"Port to listen on" default:"8080"`
}

func main() {
	//CreateTask a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		db, err := database.InitDB(options.DBPath)
		if err != nil {
			panic(err)
		}

		// CreateTask a new router & API
		router := chi.NewMux()
		api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

		entryhandlers.ListEntriesHandler(api, db)
		entryhandlers.UpdateEntryHandler(api, db)
		entryhandlers.DeleteEntryHandler(api, db)
		entryhandlers.SetTitleHandler(api, db)
		entryhandlers.SetDescriptionHandler(api, db)
		entryhandlers.ChangeRankHandler(api, db)
		entryhandlers.WeeklySummaryHandler(api, db)

		taskhandlers.CreateTaskHandler(api, db)
		taskhandlers.GetTaskSummaryHandler(api, db)
		taskhandlers.SnoozeTaskHandler(api, db)
		taskhandlers.SetTaskDoneHandler(api, db)
		taskhandlers.SetTaskUpdateHandler(api, db)
		taskhandlers.ImportPastTasksHandler(api, db)
		taskhandlers.MigrateTaskToMonthlyLogHandler(api, db)
		taskhandlers.MigrateTaskToDailyLogHandler(api, db)

		notehandlers.CreateNoteHandler(api, db)

		recurringtaskhandlers.CreateHandler(api, db)
		recurringtaskhandlers.ListHandler(api, db)
		recurringtaskhandlers.DeleteHandler(api, db)

		staticFilesHandler(router)

		// Tell the CLI how to start your router.
		hooks.OnStart(func() {
			log.Printf("Listening on port %d\n", options.Port)
			err := http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
			if err != nil {
				log.Fatalf("Error listening on port %d: %v", options.Port, err)
			}
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}

func staticFilesHandler(router *chi.Mux) {
	newRoot, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatal(err)
	}
	hfs := http.FS(newRoot)
	fserver := http.FileServer(hfs)

	router.Get("/*", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			if _, err := hfs.Open(strings.TrimPrefix(req.URL.Path, "/")); err != nil {
				http.Redirect(w, req, "/", http.StatusSeeOther)
				return
			}
		}
		fserver.ServeHTTP(w, req)
	})
}
