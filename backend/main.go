package main

import (
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"programmerjournal-backend/handlershuma"
	"programmerjournal-backend/taskrepository"

	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
)

// Options for the CLI. Pass `--port` or set the `SERVICE_PORT` env var.
type Options struct {
	DBPath string `help:"Path to the database file" default:"./foo.db"`
	Port   int    `help:"Port to listen on" default:"8080"`
}

// GreetingOutput represents the greeting operation response.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func main() {
	//var dbPath = "./foo.db"
	////var port = "8080"
	//if len(os.Args) >= 2 {
	//	dbPath = os.Args[1]
	//}
	//if len(os.Args) >= 3 {
	//	port = os.Args[2]
	//}

	//db, err := taskrepository.InitDB(dbPath)
	//if err != nil {
	//	panic(err)
	//}
	//dbRepo, err := taskrepository.NewRepository(db)
	//if err != nil {
	//	panic(err)
	//}

	//r := handlers.NewRouter(dbRepo)
	//
	//log.Printf("Starting server on port %s", port)
	//listenStr := ":" + port
	//log.Fatal(http.ListenAndServe(listenStr, r))

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

		handlershuma.ListTasks(api, dbRepo)
		handlershuma.CreateTask(api, dbRepo)
		handlershuma.UpdateTask(api, dbRepo)
		handlershuma.GetTaskSummary(api, dbRepo)
		handlershuma.DeleteTask(api, dbRepo)
		handlershuma.SnoozeTask(api, dbRepo)
		handlershuma.SetTaskDone(api, dbRepo)
		handlershuma.SetTaskTitle(api, dbRepo)
		handlershuma.SetTaskUpdate(api, dbRepo)
		handlershuma.SetTaskDescription(api, dbRepo)
		handlershuma.ChangeTaskRank(api, dbRepo)
		handlershuma.ImportPastTasks(api, dbRepo)

		// Tell the CLI how to start your router.
		hooks.OnStart(func() {
			http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
		})
	})

	// Run the CLI. When passed no commands, it starts the server.
	cli.Run()
}
