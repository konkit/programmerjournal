package recurringtask

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/recurringtask"
)

type ListEntriesInput struct {
}

type ListEntriesOutput struct {
	Body   []recurringtask.RecurringTask
	Status int
}

func ListHandler(api huma.API, rts *database.RecurringTaskService) {
	op := huma.Operation{
		OperationID: "List",
		Method:      http.MethodGet,
		Path:        "/api/recurring/list",
		Tags:        []string{"RecurringTaskService"},
	}
	huma.Register(api, op, func(ctx context.Context, input *ListEntriesInput) (*ListEntriesOutput, error) {
		resp := &ListEntriesOutput{}

		rt, err := rts.FindAll()
		if err != nil {
			return nil, err
		}

		resp.Body = rt
		resp.Status = http.StatusOK
		return resp, nil
	})
}

//func List(db *gorm.DB) ([]recurringtask.RecurringTaskService, error) {
//	var entriesFromDB []recurringtask.RecurringTaskService
//	err := db.Model(recurringtask.RecurringTaskService{}).
//		Find(&entriesFromDB).
//		Error
//
//	return entriesFromDB, err
//}
