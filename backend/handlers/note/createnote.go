package note

import (
	"context"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"net/http"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
)

type CreateNoteInput struct {
	Body struct {
		Title       string `json:"title"`
		CreatedDate string `json:"createdDate"`
	}
}

type CreateNoteResponse struct {
	Status int
}

func CreateNoteHandler(api huma.API, db *gorm.DB) {
	op := huma.Operation{
		OperationID: "CreateNote",
		Method:      http.MethodPost,
		Path:        "/api/notes/create",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *CreateNoteInput) (*CreateNoteResponse, error) {
		resp := &CreateNoteResponse{}

		dateType := date.GetDateType(input.Body.CreatedDate)
		if dateType == date.DateTypeUnrecognized {
			resp.Status = http.StatusBadRequest
			return nil, fmt.Errorf("createdDate in unrecognized date format: %s", input.Body.CreatedDate)
		}

		newTask := entry.Entry{
			Title:       input.Body.Title,
			CreatedDate: date.DateString(input.Body.CreatedDate),
		}

		err := CreateNote(db, newTask)
		if err != nil {
			return nil, err
		}

		resp.Status = http.StatusCreated
		return resp, nil
	})
}

func CreateNote(db *gorm.DB, newTask entry.Entry) error {
	count := fetchNextRank(db, newTask.CreatedDate)

	newTask.TaskID = uuid.NewString()
	newTask.Status = entry.StatusNote
	newTask.TaskUpdate = ""
	newTask.Rank = count

	return db.Create(&newTask).Error
}

func fetchNextRank(db *gorm.DB, date date.DateString) int {
	var count int64
	db.Model(entry.Entry{}).
		Where("created_date = ?", date).
		Where("rank >= 0").
		Count(&count)
	return int(count)
}
