package entry

import (
	"context"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"
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

	entries, err := es.FindEntriesByDate(entryToMove.CreatedDate)
	if err != nil {
		return err
	}

	others := make([]*entry.Entry, 0, len(entries)-1)
	originalRanks := make(map[uint]int)

	for i := range entries {
		e := &entries[i]
		originalRanks[e.ID] = e.Rank
		if e.ID != entryToMove.ID {
			others = append(others, e)
		}
	}

	positives := 0
	for _, e := range others {
		if e.Rank >= 0 {
			positives++
		}
	}
	if newIndex > positives {
		newIndex = positives
	}

	oldRank := entryToMove.Rank

	// "Remove" Step: Shift ranks down to close the gap
	for _, e := range others {
		if e.Rank > oldRank {
			e.Rank--
		}
	}

	// "Insert" Step: Shift ranks up to make space for the new entry
	// Build map for quick lookup of current state
	rankMap := make(map[int]*entry.Entry)
	for _, e := range others {
		rankMap[e.Rank] = e
	}

	// Find contiguous chain starting at newIndex
	curr := newIndex
	var chain []*entry.Entry
	for {
		if e, ok := rankMap[curr]; ok {
			chain = append(chain, e)
			curr++
		} else {
			break
		}
	}

	// Shift chain
	for _, e := range chain {
		e.Rank++
	}

	entryToMove.Rank = newIndex

	// Save changes
	if entryToMove.Rank != originalRanks[entryToMove.ID] {
		if err := es.UpdateEntry(&entryToMove); err != nil {
			return err
		}
	}

	for _, e := range others {
		if e.Rank != originalRanks[e.ID] {
			if err := es.UpdateEntry(e); err != nil {
				return err
			}
		}
	}

	return nil
}
