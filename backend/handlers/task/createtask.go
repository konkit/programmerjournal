package task

import (
	"context"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"net/http"
	"programmerjournal-backend/database"
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

func CreateTaskHandler(api huma.API, es *database.EntryService) {
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

		title := input.Body.Title
		createdDate := input.Body.CreatedDate
		err := CreateTask(es, title, createdDate)
		if err != nil {
			return nil, err
		}

		resp.Status = http.StatusCreated
		return resp, nil
	})
}

func CreateTask(es *database.EntryService, title string, createdDate string) error {
	newTask := entry.Entry{
		TaskID:      uuid.NewString(),
		Status:      entry.StatusTaskCreated,
		TaskUpdate:  "",
		Title:       title,
		CreatedDate: date.DateString(createdDate),
	}

	return es.InsertEntry(&newTask)
}
