package recurringtaskhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/recurringtask"
)

type DeleteRecurringTaskInput struct {
	ID uint64 `path:"id" example:"123" doc:"ID of the task entry"`
}

type DeleteRecurringTaskOutput struct {
	Status int
}

func Delete(api huma.API, service *recurringtask.Service) {
	op := huma.Operation{
		OperationID: "Delete",
		Method:      http.MethodDelete,
		Path:        "/api/recurring/{id}/delete",
		Tags:        []string{"RecurringTask"},
	}
	huma.Register(api, op, func(ctx context.Context, input *DeleteRecurringTaskInput) (*DeleteRecurringTaskOutput, error) {
		resp := &DeleteRecurringTaskOutput{}
		err := service.Delete(input.ID)
		if err != nil {
			return nil, err
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}
