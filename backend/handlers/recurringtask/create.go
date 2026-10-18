package recurringtask

import (
	"context"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/recurringtask"
	"time"

	"github.com/danielgtaylor/huma/v2"
)

type CreateRecurringTaskInput struct {
	Body CreateRecurringTaskInputBody
}

type CreateRecurringTaskInputBody struct {
	TaskTitle       string `json:"taskTitle"`
	TaskDescription string `json:"taskDescription"`
	FreqByWeekDay   string `json:"freqByWeekDay"`
}

type CreateRecurringTaskOutput struct {
	Status int
}

func CreateHandler(api huma.API, rts *database.RecurringTaskService) {
	op := huma.Operation{
		OperationID: "Create",
		Method:      http.MethodPost,
		Path:        "/api/recurring/create",
		Tags:        []string{"RecurringTask"},
	}
	huma.Register(api, op, func(ctx context.Context, input *CreateRecurringTaskInput) (*CreateRecurringTaskOutput, error) {
		resp := &CreateRecurringTaskOutput{}

		rTask := recurringtask.RecurringTask{
			TaskTitle:       input.Body.TaskTitle,
			TaskDescription: input.Body.TaskDescription,
			FreqByWeekDay:   input.Body.FreqByWeekDay,
			// Start generating from today, without backfilling days before the task existed
			GeneratedUntil: date.DateString(time.Now().AddDate(0, 0, -1).Format("2006-01-02")),
		}

		err := rts.Create(rTask)
		if err != nil {
			slog.Error("Error in CreateHandler", "error", err)
			return nil, huma.Error500InternalServerError(err.Error())
		}

		resp.Status = http.StatusCreated
		return resp, nil
	})
}

//func Create(db *gorm.DB, newRTask recurringtask.RecurringTaskService) error {
//	return db.Create(&newRTask).Error
//}
