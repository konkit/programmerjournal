package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/date"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
)

type CreateTaskInput struct {
	Body struct {
		Title       string `json:"title"`
		CreatedDate string `json:"createdDate"`
	}
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
		newTask := task.Task{
			Title:       input.Body.Title,
			CreatedDate: date.Parse(input.Body.CreatedDate),
		}

		err := taskRepo.Create(newTask)
		if err != nil {
			return nil, err
		}

		resp := &CreateTaskResponse{}
		resp.Status = http.StatusCreated
		return resp, nil
	})
}
