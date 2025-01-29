package recurringtask

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"net/http"
	"programmerjournal-backend/model/recurringtask"
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

func CreateHandler(api huma.API, db *gorm.DB) {
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
		}

		err := Create(db, rTask)
		if err != nil {
			return nil, err
		}

		resp.Status = http.StatusCreated
		return resp, nil
	})
}

func Create(db *gorm.DB, newRTask recurringtask.RecurringTask) error {
	return db.Create(&newRTask).Error
}
