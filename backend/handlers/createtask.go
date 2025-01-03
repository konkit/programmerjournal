package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
)

type CreateTaskInput struct {
	Body task.Task
}

type CreateTaskResponse struct {
	Status int
}

func CreateTask(api huma.API, taskRepo *taskrepository.DBRepository) {
	huma.Post(api, "/api/tasks/create", func(ctx context.Context, input *CreateTaskInput) (*CreateTaskResponse, error) {
		newTask := input.Body

		resp := &CreateTaskResponse{}
		err := taskRepo.Create(newTask)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusCreated
		return resp, nil
	})
}
