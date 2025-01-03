package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
)

type UpdateTaskInput struct {
	Body task.Task
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
}

type UpdateTaskResponse struct {
	Status int
}

func UpdateTask(api huma.API, taskRepo *taskrepository.DBRepository) {
	huma.Put(api, "/api/tasks/{id}/update", func(ctx context.Context, input *UpdateTaskInput) (*UpdateTaskResponse, error) {
		resp := &UpdateTaskResponse{}
		err := taskRepo.Update(input.ID, input.Body)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
