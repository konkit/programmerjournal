package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"programmerjournal-backend/database"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	_ "github.com/danielgtaylor/huma/v2/formats/cbor"
	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/go-chi/chi/v5"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/gorm"

	entryhandlers "programmerjournal-backend/handlers/entry"
	notehandlers "programmerjournal-backend/handlers/note"
	recurringtaskhandlers "programmerjournal-backend/handlers/recurringtask"
	taghandlers "programmerjournal-backend/handlers/tags"
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
	cli := humacli.New(func(hooks humacli.Hooks, options *Options) {
		db, err := database.InitDB(options.DBPath)
		if err != nil {
			panic(err)
		}

		router := chi.NewMux()
		staticFilesHandler(router)

		api := humachi.New(router, huma.DefaultConfig("My API", "1.0.0"))
		registerHandlers(api, db)

		hooks.OnStart(func() {
			log.Printf("Listening on port %d\n", options.Port)
			err := http.ListenAndServe(fmt.Sprintf(":%d", options.Port), router)
			if err != nil {
				log.Fatalf("Error listening on port %d: %v", options.Port, err)
			}
		})
	})

	cli.Run()
}

func registerHandlers(api huma.API, db *gorm.DB) {
	es := database.NewEntryService(db)
	rts := database.NewRecurringTaskService(db)

	entryhandlers.ListEntriesHandler(api, es)
	entryhandlers.DeleteEntryHandler(api, es)
	entryhandlers.SetTitleHandler(api, es)
	entryhandlers.SetDescriptionHandler(api, es)
	entryhandlers.ChangeRankHandler(api, es)
	entryhandlers.WeeklySummaryHandler(api, es)
	entryhandlers.WeeklyUpdatesHandler(api, es)

	taskhandlers.CreateTaskHandler(api, es)
	taskhandlers.GetTaskSummaryHandler(api, es)
	taskhandlers.SnoozeTaskHandler(api, es)
	taskhandlers.SetTaskDoneHandler(api, es)
	taskhandlers.CancelTaskHandler(api, es)
	taskhandlers.SetTaskUpdateHandler(api, es)
	taskhandlers.ImportPastTasksHandler(api, rts, es)
	taskhandlers.CountPastTasks(api, rts, es)
	taskhandlers.MigrateTaskToMonthlyLogHandler(api, es)
	taskhandlers.MigrateTaskToDailyLogHandler(api, es)

	notehandlers.CreateNoteHandler(api, es)

	recurringtaskhandlers.CreateHandler(api, rts)
	recurringtaskhandlers.ListHandler(api, rts)
	recurringtaskhandlers.DeleteHandler(api, rts)

	taghandlers.ListTagsHandler(api, db)
	taghandlers.GetTagsHandler(api, db)
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
