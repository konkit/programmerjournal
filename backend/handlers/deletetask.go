package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/taskrepository"
)

type DeleteTaskInput struct {
	ID uint64 `path:"id" example:"123" doc:"ID of the task entry"`
}

type DeleteTaskResponse struct {
	Status int
}

func DeleteTask(api huma.API, taskRepo *taskrepository.DBRepository) {
	op := huma.Operation{
		OperationID: "DeleteTask",
		Method:      http.MethodDelete,
		Path:        "/api/tasks/{id}/delete",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *DeleteTaskInput) (*DeleteTaskResponse, error) {
		resp := &DeleteTaskResponse{}
		err := taskRepo.Delete(input.ID)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
