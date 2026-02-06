package entry

import (
	"context"
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

	entriesWithoutElementToBeMoved := filterEntriesWithoutElementToBeMoved(entriesFromDB, oldIndex)

	elementAdded := false
	var currentRank int
	if len(entriesWithoutElementToBeMoved) > 0 {
		currentRank = min(entriesWithoutElementToBeMoved[0].Rank, newIndex, 0)
		entriesIter := 0

		for entriesIter < len(entriesWithoutElementToBeMoved) {
			if currentRank == newIndex {
				err = saveIfRankChanged(es, e, currentRank)
				if err != nil {
					return err
				}
				elementAdded = true
			} else {
				entryToAdd := entriesWithoutElementToBeMoved[entriesIter]
				if currentRank >= 0 || entryToAdd.Rank < 0 {
					entriesIter++
					err = saveIfRankChanged(es, entryToAdd, currentRank)
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
		err := appendElementIfNotYetSet(es, e, currentRank, newIndex)
		if err != nil {
			return err
		}
	}

	return nil
}

func filterEntriesWithoutElementToBeMoved(entriesFromDB []entry.Entry, oldIndex int) []entry.Entry {
	var entriesWithoutMovedElement []entry.Entry
	elementFound := false
	for _, en := range entriesFromDB {
		if elementFound != true && en.Rank == oldIndex {
			elementFound = true
			continue
		}

		entriesWithoutMovedElement = append(entriesWithoutMovedElement, en)
	}
	return entriesWithoutMovedElement
}

func appendElementIfNotYetSet(es *database.EntryService, e entry.Entry, currentRank int, newIndex int) error {
	// Set the initial index as zero if moving from the priority list back to the normal one
	if currentRank < 0 && newIndex >= 0 {
		currentRank = 0
	}
	return saveIfRankChanged(es, e, currentRank)
}

func saveIfRankChanged(es *database.EntryService, e entry.Entry, currentRank int) error {
	// Do not save if the rank is already set to currentRank
	if e.Rank != currentRank {
		e.Rank = currentRank
		return es.UpdateEntry(&e)
	}
	return nil
}
