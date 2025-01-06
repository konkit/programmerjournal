package handlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/date"
	"programmerjournal-backend/task"
	"programmerjournal-backend/taskrepository"
)

type CreateNoteInput struct {
	Body struct {
		Title       string `json:"title"`
		CreatedDate string `json:"createdDate"`
	}
}

type CreateNoteResponse struct {
	Status int
}

func CreateNote(api huma.API, taskRepo *taskrepository.DBRepository) {
	op := huma.Operation{
		OperationID: "CreateNote",
		Method:      http.MethodPost,
		Path:        "/api/notes/create",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *CreateNoteInput) (*CreateNoteResponse, error) {
		newTask := task.Task{
			Title:       input.Body.Title,
			CreatedDate: date.Parse(input.Body.CreatedDate),
		}

		err := taskRepo.CreateNote(newTask)
		if err != nil {
			return nil, err
		}

		resp := &CreateNoteResponse{}
		resp.Status = http.StatusCreated
		return resp, nil
	})
}
