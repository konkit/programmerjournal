package task

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/entry"
)

type GetTaskSummaryInput struct {
	ID uint64 `path:"id" example:"123" doc:"ID of the task entry"`
}

type GetTaskSummaryResponse struct {
	Status int
	Body   *entry.TaskSummary
}

func GetTaskSummaryHandler(api huma.API, es *database.Entry) {
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
			return nil, err
		}

		resp.Body = summary
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func GetTaskSummary(es *database.Entry, id uint64) (*entry.TaskSummary, error) {
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
