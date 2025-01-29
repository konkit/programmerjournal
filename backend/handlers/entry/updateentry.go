package entry

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"net/http"
	"programmerjournal-backend/model/entry"
)

type UpdateEntryInput struct {
	Body entry.Entry
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
}

type UpdateEntryResponse struct {
	Status int
}

func UpdateEntryHandler(api huma.API, db *gorm.DB) {
	op := huma.Operation{
		OperationID: "UpdateEntry",
		Method:      http.MethodPut,
		Path:        "/api/tasks/{id}/update",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *UpdateEntryInput) (*UpdateEntryResponse, error) {
		resp := &UpdateEntryResponse{}
		err := Update(db, input.ID, input.Body)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func Update(db *gorm.DB, entryID uint64, updatedTask entry.Entry) error {
	updatedTask.ID = uint(entryID)
	return db.Save(updatedTask).Error
}
