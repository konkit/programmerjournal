package note

import (
	"context"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/entry"

	"github.com/danielgtaylor/huma/v2"
)

type ListNotesResponse struct {
	Body struct {
		Notes []entry.Entry `json:"notes"`
	}
}

func ListNotesHandler(api huma.API, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "ListNotes",
		Method:      http.MethodGet,
		Path:        "/api/notes",
		Tags:        []string{"Note"},
	}
	huma.Register(api, op, func(ctx context.Context, input *struct{}) (*ListNotesResponse, error) {
		resp := &ListNotesResponse{}

		notes, err := es.FindAllNotes()
		if err != nil {
			slog.Error("Error in ListNotesHandler", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp.Body.Notes = notes
		return resp, nil
	})
}
