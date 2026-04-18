package entry

import (
	"context"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"

	"github.com/danielgtaylor/huma/v2"
)

type DeleteEntryInput struct {
	ID uint `path:"id" example:"123" doc:"ID of the task entry"`
}

type DeleteEntryResponse struct {
	Status int
}

func DeleteEntryHandler(api huma.API, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "DeleteEntry",
		Method:      http.MethodDelete,
		Path:        "/api/entries/{id}/delete",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *DeleteEntryInput) (*DeleteEntryResponse, error) {
		resp := &DeleteEntryResponse{}
		err := DeleteEntry(es, input)
		if err != nil {
			slog.Error("Error in DeleteEntryHandler", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func DeleteEntry(es *database.EntryService, input *DeleteEntryInput) error {
	e, err := es.GetEntryByID(input.ID)
	if err != nil {
		return err
	}

	err = es.DeleteEntry(input.ID)
	if err != nil {
		return err
	}

	// Adjust Ranks of other items
	entriesFromDB, err := es.FindEntriesByDate(e.CreatedDate)
	if err != nil {
		return err
	}

	for i, en := range entriesFromDB {
		en.Rank = i
		err = es.UpdateEntry(&en)
		if err != nil {
			return err
		}
	}

	return nil
}
