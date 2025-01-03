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
	op := huma.Operation{
		OperationID: "CreateTask",
		Method:      http.MethodPost,
		Path:        "/api/tasks/create",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *CreateTaskInput) (*CreateTaskResponse, error) {
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
