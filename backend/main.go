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
	"programmerjournal-backend/handlers/entryhandlers"
	"programmerjournal-backend/handlers/notehandlers"
	"programmerjournal-backend/handlers/recurringtaskhandlers"
	"programmerjournal-backend/handlers/taskhandlers"
	"programmerjournal-backend/model/database"
	"programmerjournal-backend/model/entry"
	"programmerjournal-backend/model/recurringtask"
	"strings"
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
		entryService := entry.NewService(db)
		recurringTaskService := recurringtask.NewService(db)

		// CreateTask a new router & API
		router := chi.NewMux()
		api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

		entryhandlers.ListEntriesHandler(api, entryService)
		entryhandlers.UpdateEntry(api, entryService)
		entryhandlers.DeleteEntryHandler(api, entryService)
		entryhandlers.SetTitle(api, entryService)
		entryhandlers.SetDescription(api, entryService)
		entryhandlers.ChangeRankHandler(api, entryService)
		entryhandlers.WeeklySummary(api, entryService)

		taskhandlers.CreateTask(api, entryService)
		taskhandlers.GetTaskSummary(api, entryService)
		taskhandlers.SnoozeTask(api, entryService)
		taskhandlers.SetTaskDone(api, entryService)
		taskhandlers.SetTaskUpdate(api, entryService)
		taskhandlers.ImportPastTasks(api, entryService)
		taskhandlers.MigrateTaskToMonthlyLog(api, entryService)
		taskhandlers.MigrateTaskToDailyLog(api, entryService)

		notehandlers.CreateNote(api, entryService)

		recurringtaskhandlers.Create(api, recurringTaskService)
		recurringtaskhandlers.List(api, recurringTaskService)
		recurringtaskhandlers.Delete(api, recurringTaskService)

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
