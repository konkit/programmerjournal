package recurringtask

import (
	"context"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"

	"github.com/danielgtaylor/huma/v2"
)

type DeleteRecurringTaskInput struct {
	ID uint64 `path:"id" example:"123" doc:"ID of the task entry"`
}

type DeleteRecurringTaskOutput struct {
	Status int
}

func DeleteHandler(api huma.API, rts *database.RecurringTaskService) {
	op := huma.Operation{
		OperationID: "Delete",
		Method:      http.MethodDelete,
		Path:        "/api/recurring/{id}/delete",
		Tags:        []string{"RecurringTask"},
	}
	huma.Register(api, op, func(ctx context.Context, input *DeleteRecurringTaskInput) (*DeleteRecurringTaskOutput, error) {
		resp := &DeleteRecurringTaskOutput{}
		err := rts.Delete(input.ID)
		if err != nil {
			slog.Error("Error in DeleteHandler", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
		}
		resp.Status = http.StatusOK
		return resp, nil
	})
}

//func Delete(db *gorm.DB, id uint64) error {
//	return db.Delete(&recurringtask.RecurringTaskService{}, id).Error
//}
