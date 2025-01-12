package notehandlers

import (
	"context"
	"fmt"
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

		dateType := date.GetDateType(input.Body.CreatedDate)
		if dateType == date.DateTypeUnrecognized {
			resp.Status = http.StatusBadRequest
			return nil, fmt.Errorf("createdDate in unrecognized date format: %s", input.Body.CreatedDate)
		}

		newTask := entry.Entry{
			Title:       input.Body.Title,
			CreatedDate: date.DateString(input.Body.CreatedDate),
		}

		err := taskRepo.CreateNote(newTask)
		if err != nil {
			return nil, err
		}

		resp.Status = http.StatusCreated
		return resp, nil
	})
}
