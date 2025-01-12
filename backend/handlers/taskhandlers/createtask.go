package taskhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
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

func CreateTask(api huma.API, taskRepo *entry.Service) {
	op := huma.Operation{
		OperationID: "CreateTask",
		Method:      http.MethodPost,
		Path:        "/api/tasks/create",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *CreateTaskInput) (*CreateTaskResponse, error) {
		resp := &CreateTaskResponse{}
		dateString, err := date.ParseDateString(input.Body.CreatedDate)
		if err != nil {
			resp.Status = http.StatusBadRequest
			return nil, err
		}
		newTask := entry.Entry{
			Title:       input.Body.Title,
			CreatedDate: dateString,
		}

		err = taskRepo.CreateTask(newTask)
		if err != nil {
			return nil, err
		}

		resp.Status = http.StatusCreated
		return resp, nil
	})
}
