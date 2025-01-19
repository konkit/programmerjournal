package recurringtaskhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
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

func Create(api huma.API, rTaskService *recurringtask.Service) {
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

		err := rTaskService.Create(rTask)
		if err != nil {
			return nil, err
		}

		resp.Status = http.StatusCreated
		return resp, nil
	})
}
