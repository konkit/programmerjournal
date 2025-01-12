package taskhandlers

import (
	"context"
	"fmt"
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
		dateType := date.GetDateType(input.Body.CreatedDate)
		if dateType == date.DateTypeUnrecognized {
			resp.Status = http.StatusBadRequest
			return nil, fmt.Errorf("createdDate in unrecognized date format: %s", input.Body.CreatedDate)
		}

		newTask := entry.Entry{
			Title:       input.Body.Title,
			CreatedDate: date.DateString(input.Body.CreatedDate),
		}

		err := taskRepo.CreateTask(newTask)
		if err != nil {
			return nil, err
		}

		resp.Status = http.StatusCreated
		return resp, nil
	})
}
