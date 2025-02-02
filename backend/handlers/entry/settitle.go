package entry

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"net/http"
	"programmerjournal-backend/handlers/utils"
)

type SetTitleInput struct {
	ID   uint64 `path:"id" example:"123" doc:"ID of the task entry"`
	Body struct {
		//Date string `json:"date" doc:"Date when the task should be snoozed"`
		Title string `json:"title" doc:"New title of the task"`
	}
}

type SetTaskTitleResponse struct {
	Status int
}

func SetTitleHandler(api huma.API, db *gorm.DB) {
	op := huma.Operation{
		OperationID: "SetTitle",
		Method:      http.MethodPatch,
		Path:        "/api/entries/{id}/setTitle",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *SetTitleInput) (*SetTaskTitleResponse, error) {
		resp := &SetTaskTitleResponse{}
		err := SetTitle(db, input.ID, input.Body.Title)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func SetTitle(db *gorm.DB, entryID uint64, title string) error {
	t, err := utils.GetEntryByID(db, entryID)
	if err != nil {
		return err
	}

	t.Title = title
	db.Save(&t)

	return nil
}
