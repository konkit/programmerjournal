package main

import (
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"programmerjournal-backend/handlers"
	"programmerjournal-backend/taskrepository"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

// Options for the CLI. Pass `--port` or set the `SERVICE_PORT` env var.
type Options struct {
	DBPath string `help:"Path to the database file" default:"./foo.db"`
	Port   int    `help:"Port to listen on" default:"8080"`
}

func main() {
	//Create a CLI app which takes a port option.
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		db, err := taskrepository.InitDB(options.DBPath)
		if err != nil {
			panic(err)
		}
		dbRepo, err := taskrepository.NewRepository(db)
		if err != nil {
			panic(err)
		}

		// Create a new router & API
		router := chi.NewMux()
		api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))

		handlers.ListTasks(api, dbRepo)
		handlers.CreateTask(api, dbRepo)
		handlers.UpdateTask(api, dbRepo)
		handlers.GetTaskSummary(api, dbRepo)
		handlers.DeleteTask(api, dbRepo)
		handlers.SnoozeTask(api, dbRepo)
		handlers.SetTaskDone(api, dbRepo)
		handlers.SetTaskTitle(api, dbRepo)
		handlers.SetTaskUpdate(api, dbRepo)
		handlers.SetTaskDescription(api, dbRepo)
		handlers.ChangeTaskRank(api, dbRepo)
		handlers.ImportPastTasks(api, dbRepo)
		handlers.MigrateToMonthlyLog(api, dbRepo)

		// Tell the CLI how to start your router.
		hooks.OnStart(func() {
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
