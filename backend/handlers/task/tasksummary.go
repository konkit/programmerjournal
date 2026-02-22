package task

import (
	"context"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/entry"

	"github.com/danielgtaylor/huma/v2"
)

type GetTaskSummaryInput struct {
	ID uint `path:"id" example:"123" doc:"ID of the task entry"`
}

type GetTaskSummaryResponse struct {
	Status int
	Body   *entry.TaskSummary
}

func GetTaskSummaryHandler(api huma.API, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "GetTaskSummary",
		Method:      http.MethodGet,
		Path:        "/api/tasks/{id}/summary",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *GetTaskSummaryInput) (*GetTaskSummaryResponse, error) {
		resp := &GetTaskSummaryResponse{}
		summary, err := GetTaskSummary(es, input.ID)
		if err != nil {
			slog.Error("Error in GetTaskSummaryHandler", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp.Body = summary
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func GetTaskSummary(es *database.EntryService, id uint) (*entry.TaskSummary, error) {
	e, err := es.GetEntryByID(id)
	if err != nil {
		return nil, err
	}

	tasksFromDB, err := es.FindTasksByTaskID(e.TaskID)
	if err != nil {
		return nil, err
	}

	var updates []entry.TaskUpdate
	for _, tt := range tasksFromDB {
		update := entry.TaskUpdate{
			Date:   tt.CreatedDate,
			Update: tt.TaskUpdate,
			Status: tt.Status,
		}
		updates = append(updates, update)
	}

	ts := &entry.TaskSummary{
		TaskEntry: e,
		Updates:   updates,
	}

	return ts, nil
}
