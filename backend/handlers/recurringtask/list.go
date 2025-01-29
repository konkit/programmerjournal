package recurringtask

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"net/http"
	"programmerjournal-backend/model/recurringtask"
)

type ListEntriesInput struct {
}

type ListEntriesOutput struct {
	Body   []recurringtask.RecurringTask
	Status int
}

func ListHandler(api huma.API, db *gorm.DB) {
	op := huma.Operation{
		OperationID: "List",
		Method:      http.MethodGet,
		Path:        "/api/recurring/list",
		Tags:        []string{"RecurringTask"},
	}
	huma.Register(api, op, func(ctx context.Context, input *ListEntriesInput) (*ListEntriesOutput, error) {
		resp := &ListEntriesOutput{}

		rt, err := List(db)
		if err != nil {
			return nil, err
		}

		resp.Body = rt
		resp.Status = http.StatusOK
		return resp, nil
	})
}

func List(db *gorm.DB) ([]recurringtask.RecurringTask, error) {
	var entriesFromDB []recurringtask.RecurringTask
	err := db.Model(recurringtask.RecurringTask{}).
		Find(&entriesFromDB).
		Error

	return entriesFromDB, err
}
