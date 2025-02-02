package entry

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/entry"
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
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func ChangeRank(es *database.EntryService, entryID uint, newIndex int) error {
	e, err := es.GetEntryByID(entryID)
	if err != nil {
		return err
	}
	oldIndex := e.Rank

	entriesFromDB, err := es.FindEntriesByDate(e.CreatedDate)
	if err != nil {
		return err
	}

	var entriesWithoutMovedElement []entry.Entry
	elementFound := false
	for _, en := range entriesFromDB {
		if en.Rank != oldIndex {
			entriesWithoutMovedElement = append(entriesWithoutMovedElement, en)
		} else if elementFound == true {
			entriesWithoutMovedElement = append(entriesWithoutMovedElement, en)
		} else {
			elementFound = true
		}
	}

	elementAdded := false
	var currentRank int
	if len(entriesWithoutMovedElement) > 0 {
		currentRank = min(entriesWithoutMovedElement[0].Rank, newIndex)
		entriesIter := 0

		for entriesIter < len(entriesWithoutMovedElement) {
			if currentRank == newIndex {
				err = saveWithNewRank(es, e, currentRank)
				if err != nil {
					return err
				}
				elementAdded = true
			} else {
				entryToAdd := entriesWithoutMovedElement[entriesIter]
				if currentRank >= 0 || entryToAdd.Rank < 0 {
					entriesIter++
					err = saveWithNewRank(es, entryToAdd, currentRank)
					if err != nil {
						return err
					}
				}
			}

			currentRank++
		}
	} else {
		currentRank = min(newIndex, 0)
	}

	if !elementAdded {
		if currentRank < 0 && newIndex >= 0 {
			currentRank = 0
		}
		err = saveWithNewRank(es, e, currentRank)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveWithNewRank(es *database.EntryService, e entry.Entry, currentRank int) error {
	// Do not save if the rank is already set to currentRank
	if e.Rank != currentRank {
		e.Rank = currentRank
		return es.UpdateEntry(&e)
	}
	return nil
}
