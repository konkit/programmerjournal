package entry

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"net/http"
	"programmerjournal-backend/model/entry"
)

type ChangeRankInput struct {
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		NewRank int `json:"newRank"`
	}
}

type ChangeRankResponse struct {
	Status int
}

func ChangeRankHandler(api huma.API, db *gorm.DB) {
	op := huma.Operation{
		OperationID: "ChangeRank",
		Method:      http.MethodPatch,
		Path:        "/api/entries/{id}/changeRank",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *ChangeRankInput) (*ChangeRankResponse, error) {
		resp := &ChangeRankResponse{}
		err := ChangeRank(db, input.ID, input.Body.NewRank)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func ChangeRank(db *gorm.DB, entryID uint64, newIndex int) error {
	e, err := getEntryByID(db, entryID)
	if err != nil {
		return err
	}
	oldIndex := e.Rank

	var entriesFromDB []entry.Entry
	err = db.Model(entry.Entry{}).
		Order("rank").
		Where("created_date = ?", e.CreatedDate).
		Find(&entriesFromDB).
		Error

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
				err = saveWithNewRank(db, e, currentRank)
				if err != nil {
					return err
				}
				elementAdded = true
			} else {
				entryToAdd := entriesWithoutMovedElement[entriesIter]
				if currentRank >= 0 || entryToAdd.Rank < 0 {
					entriesIter++
					err = saveWithNewRank(db, entryToAdd, currentRank)
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
		err = saveWithNewRank(db, e, currentRank)
		if err != nil {
			return err
		}
	}

	return nil
}

func getEntryByID(db *gorm.DB, entryID uint64) (entry.Entry, error) {
	t := entry.Entry{ID: uint(entryID)}
	if err := db.First(&t).Error; err != nil {
		return t, err
	}
	return t, nil
}

func saveWithNewRank(db *gorm.DB, e entry.Entry, currentRank int) error {
	// Do not save if the rank is already set to currentRank
	if e.Rank != currentRank {
		e.Rank = currentRank
		return db.Save(e).Error
	}
	return nil
}
