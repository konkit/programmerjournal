package taskhandlers

import (
	"context"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
)

type ImportPastTasksInput struct {
	Date string `path:"date"`
}

type ImportPastTasksResponse struct {
	Status int
}

func ImportPastTasks(api huma.API, taskRepo *entry.Service) {
	op := huma.Operation{
		OperationID: "ImportPastTasks",
		Method:      http.MethodPost,
		Path:        "/api/tasks/importPastTasks/{date}",
		Tags:        []string{"Task"},
	}
	huma.Register(api, op, func(ctx context.Context, input *ImportPastTasksInput) (*ImportPastTasksResponse, error) {
		resp := &ImportPastTasksResponse{}

		inputDateString := date.DateString(input.Date)

		dayDate, err := date.ParseDayDate(inputDateString)
		if err != nil {
			monthDate, err := date.ParseMonthDate(inputDateString)
			if err != nil {
				resp.Status = http.StatusBadRequest
				return resp, fmt.Errorf("invalid date format: %s", input.Date)
			} else {
				err := taskRepo.ImportPastTasksFromMonth(monthDate)
				if err != nil {
					return nil, err
				}
			}
		} else {
			err := taskRepo.ImportPastTasksFromDay(dayDate)
			if err != nil {
				return nil, err
			}
		}

		resp.Status = http.StatusOK
		return resp, nil
	})
}
