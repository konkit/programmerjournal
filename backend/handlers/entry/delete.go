package entry

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/entry"
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
			return nil, err
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

	var nonPriority []entry.Entry
	for _, en := range entriesFromDB {
		if en.Rank >= 0 {
			nonPriority = append(nonPriority, en)
		}
	}

	for i, en := range nonPriority {
		en.Rank = i
		err = es.UpdateEntry(&en)
		if err != nil {
			return err
		}
	}

	return nil
}
