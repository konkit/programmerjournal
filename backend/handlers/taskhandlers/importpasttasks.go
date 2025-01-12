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

		dateType := date.GetDateType(input.Date)

		switch dateType {
		case date.DateTypeDay:
			dayDate, err := date.ParseDayDate(date.DateString(input.Date))
			if err != nil {
				resp.Status = http.StatusBadRequest
				return nil, err
			}

			err = taskRepo.ImportPastTasksFromDay(dayDate)
			if err != nil {
				return nil, err
			}

			resp.Status = http.StatusOK
			return resp, nil
		case date.DateTypeMonth:
			monthDate, err := date.ParseMonthDate(date.DateString(input.Date))
			if err != nil {
				resp.Status = http.StatusBadRequest
				return nil, err
			}

			err = taskRepo.ImportPastTasksFromMonth(monthDate)
			if err != nil {
				return nil, err
			}

			resp.Status = http.StatusOK
			return resp, nil
		default:
			resp.Status = http.StatusBadRequest
			return nil, fmt.Errorf("unrecognized date format: %s", input.Date)
		}
	})
}
