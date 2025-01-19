package recurringtaskhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/recurringtask"
)

type ListEntriesInput struct {
}

type ListEntriesOutput struct {
	Body   []recurringtask.RecurringTask
	Status int
}

func List(api huma.API, service *recurringtask.Service) {
	op := huma.Operation{
		OperationID: "List",
		Method:      http.MethodGet,
		Path:        "/api/recurring/list",
		Tags:        []string{"RecurringTask"},
	}
	huma.Register(api, op, func(ctx context.Context, input *ListEntriesInput) (*ListEntriesOutput, error) {
		resp := &ListEntriesOutput{}

		rt, err := service.List()
		if err != nil {
			return nil, err
		}

		resp.Body = rt
		resp.Status = http.StatusOK
		return resp, nil
	})
}
