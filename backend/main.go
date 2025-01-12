package main

import (
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"programmerjournal-backend/handlers/entryhandlers"
	"programmerjournal-backend/handlers/notehandlers"
	"programmerjournal-backend/handlers/taskhandlers"
	"programmerjournal-backend/model/entry"
)

// Options for the CLI. Pass `--port` or set the `SERVICE_PORT` env var.
type Options struct {
	DBPath string `help:"Path to the database file" default:"./foo.db"`
	Port   int    `help:"Port to listen on" default:"8080"`
}

func main() {
	//CreateTask a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		db, err := entry.InitDB(options.DBPath)
		if err != nil {
			panic(err)
		}
		dbRepo, err := entry.NewService(db)
		if err != nil {
			panic(err)
		}

		// CreateTask a new router & API
		router := chi.NewMux()
		api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

		entryhandlers.ListEntries(api, dbRepo)
		entryhandlers.UpdateEntry(api, dbRepo)
		entryhandlers.DeleteEntry(api, dbRepo)
		entryhandlers.SetTitle(api, dbRepo)
		entryhandlers.SetDescription(api, dbRepo)
		entryhandlers.ChangeRank(api, dbRepo)

		taskhandlers.CreateTask(api, dbRepo)
		taskhandlers.GetTaskSummary(api, dbRepo)
		taskhandlers.SnoozeTask(api, dbRepo)
		taskhandlers.SetTaskDone(api, dbRepo)
		taskhandlers.SetTaskUpdate(api, dbRepo)
		taskhandlers.ImportPastTasks(api, dbRepo)
		taskhandlers.MigrateTaskToMonthlyLog(api, dbRepo)
		taskhandlers.WeeklyTaskSummary(api, dbRepo)

		notehandlers.CreateNote(api, dbRepo)

		// Tell the CLI how to start your router.
		hooks.OnStart(func() {
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
