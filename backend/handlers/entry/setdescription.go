package entry

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"net/http"
	"programmerjournal-backend/database"
)

type SetDescriptionInput struct {
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		//Date string `json:"date" doc:"Date when the task should be snoozed"`
		Description string `json:"description" doc:"New title of the task"`
	}
}

type SetTaskDescriptionResponse struct {
	Status int
}

func SetDescriptionHandler(api huma.API, db *gorm.DB) {
	op := huma.Operation{
		OperationID: "SetDescription",
		Method:      http.MethodPatch,
		Path:        "/api/entries/{id}/setDescription",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *SetDescriptionInput) (*SetTaskDescriptionResponse, error) {
		resp := &SetTaskDescriptionResponse{}
		err := SetDescription(db, input.ID, input.Body.Description)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func SetDescription(db *gorm.DB, entryID uint64, description string) error {
	t, err := database.GetEntryByID(db, entryID)
	if err != nil {
		return err
	}

	t.Description = description

	return database.UpdateEntry(db, &t)
}
