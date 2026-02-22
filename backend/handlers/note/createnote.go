package note

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
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

func CreateNoteHandler(api huma.API, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "CreateNote",
		Method:      http.MethodPost,
		Path:        "/api/notes/create",
		Tags:        []string{"Note"},
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

		err := CreateNote(es, newTask)
		if err != nil {
			slog.Error("Error in CreateNoteHandler", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp.Status = http.StatusCreated
		return resp, nil
	})
}

func CreateNote(es *database.EntryService, newTask entry.Entry) error {
	newTask.TaskID = uuid.NewString()
	newTask.Status = entry.StatusNote
	newTask.TaskUpdate = ""

	return es.InsertEntry(&newTask)
}
