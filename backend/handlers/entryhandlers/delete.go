package entryhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"net/http"
	"programmerjournal-backend/model/entry"
)

type DeleteEntryInput struct {
	ID uint64 `path:"id" example:"123" doc:"ID of the task entry"`
}

type DeleteEntryResponse struct {
	Status int
}

func DeleteEntryHandler(api huma.API, db *gorm.DB) {
	op := huma.Operation{
		OperationID: "DeleteEntry",
		Method:      http.MethodDelete,
		Path:        "/api/entries/{id}/delete",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *DeleteEntryInput) (*DeleteEntryResponse, error) {
		resp := &DeleteEntryResponse{}
		err := Delete(db, input.ID)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func Delete(db *gorm.DB, entryID uint64) error {
	// TODO: Handle Rank change if the task is deleted in the middle
	return db.Delete(&entry.Entry{}, entryID).Error
}
