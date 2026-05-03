package entry

import (
	"context"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"

	"github.com/danielgtaylor/huma/v2"
)

type ChangeRankInput struct {
	ID   uint `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		NewRank int `json:"newRank"`
	}
}

type ChangeRankResponse struct {
	Status int
}

func ChangeRankHandler(api huma.API, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "ChangeRank",
		Method:      http.MethodPatch,
		Path:        "/api/entries/{id}/changeRank",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *ChangeRankInput) (*ChangeRankResponse, error) {
		resp := &ChangeRankResponse{}
		err := ChangeRank(es, input.ID, input.Body.NewRank)
		if err != nil {
			slog.Error("Error in ChangeRankHandler", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func ChangeRank(es *database.EntryService, entryID uint, newIndex int) error {
	entryToMove, err := es.GetEntryByID(entryID)
	if err != nil {
		return err
	}

	if entryToMove.Status == entry.StatusTaskMigrated || entryToMove.Status == entry.StatusTaskCancelled || entryToMove.Status == entry.StatusTaskSnoozed {
		return nil
	}

	entries, err := es.FindActiveEntriesByDate(entryToMove.CreatedDate)
	if err != nil {
		return err
	}

	var others []*entry.Entry
	for i := range entries {
		e := &entries[i]
		if e.ID != entryToMove.ID {
			others = append(others, e)
		}
	}

	if newIndex > len(others) {
		newIndex = len(others)
	}
	if newIndex < 0 {
		newIndex = 0
	}

	var newOrder []*entry.Entry
	newOrder = append(newOrder, others[:newIndex]...)
	newOrder = append(newOrder, &entryToMove)
	newOrder = append(newOrder, others[newIndex:]...)

	for i, e := range newOrder {
		if e.Rank != i {
			e.Rank = i
			if err := es.UpdateEntry(e); err != nil {
				return err
			}
		}
	}

	return nil
}

func ReRankActiveTasks(es *database.EntryService, dateString date.DateString) error {
	entries, err := es.FindActiveEntriesByDate(dateString)
	if err != nil {
		return err
	}

	for i := range entries {
		e := &entries[i]
		if e.Rank != i {
			e.Rank = i
			err = es.UpdateEntry(e)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
