package notehandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
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

func CreateNote(api huma.API, taskRepo *entry.Service) {
	op := huma.Operation{
		OperationID: "CreateNote",
		Method:      http.MethodPost,
		Path:        "/api/notes/create",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *CreateNoteInput) (*CreateNoteResponse, error) {
		resp := &CreateNoteResponse{}
		dateString, err := date.ParseDateString(input.Body.CreatedDate)
		if err != nil {
			resp.Status = http.StatusBadRequest
			return nil, err
		}

		newTask := entry.Entry{
			Title:       input.Body.Title,
			CreatedDate: dateString,
		}

		err = taskRepo.CreateNote(newTask)
		if err != nil {
			return nil, err
		}

		resp.Status = http.StatusCreated
		return resp, nil
	})
}
