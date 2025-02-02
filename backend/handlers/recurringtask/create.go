package recurringtask

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/database"
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

func CreateHandler(api huma.API, rts *database.RecurringTaskService) {
	op := huma.Operation{
		OperationID: "Create",
		Method:      http.MethodPost,
		Path:        "/api/recurring/create",
		Tags:        []string{"RecurringTaskService"},
	}
	huma.Register(api, op, func(ctx context.Context, input *CreateRecurringTaskInput) (*CreateRecurringTaskOutput, error) {
		resp := &CreateRecurringTaskOutput{}

		rTask := recurringtask.RecurringTask{
			TaskTitle:       input.Body.TaskTitle,
			TaskDescription: input.Body.TaskDescription,
			FreqByWeekDay:   input.Body.FreqByWeekDay,
		}

		err := rts.Create(rTask)
		if err != nil {
			return nil, err
		}

		resp.Status = http.StatusCreated
		return resp, nil
	})
}

//func Create(db *gorm.DB, newRTask recurringtask.RecurringTaskService) error {
//	return db.Create(&newRTask).Error
//}
