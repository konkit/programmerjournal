package entryhandlers

import (
	"context"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
)

type ListEntriesInput struct {
	Date string `path:"date" example:"2024-05-05" doc:"Day to select the list from"`
}

type ListEntriesOutput struct {
	Body   []entry.Entry
	Status int
}

func ListEntries(api huma.API, taskRepo *entry.Service) {
	op := huma.Operation{
		OperationID: "ListEntries",
		Method:      http.MethodGet,
		Path:        "/api/entries/list/{date}",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *ListEntriesInput) (*ListEntriesOutput, error) {
		resp := &ListEntriesOutput{}
		var entries []entry.Entry
		var err error

		dayDate, err := date.ParseDayDate(date.DateString(input.Date))
		if err != nil {
			monthDate, err := date.ParseMonthDate(date.DateString(input.Date))
			if err != nil {
				resp.Status = http.StatusBadRequest
				return resp, fmt.Errorf("invalid date format: %s", input.Date)
			} else {
				entries, err = taskRepo.ListMonthEntries(monthDate)
				if err != nil {
					return nil, err
				}
			}
		} else {
			entries, err = taskRepo.ListDayEntries(dayDate)
			if err != nil {
				return nil, err
			}
		}

		//tasks, err := taskRepo.ListEntries(date.Parse(input.Date))
		if err != nil {
			return nil, err
		}
		resp.Body = entries
		resp.Status = http.StatusOK
		return resp, nil
	})
}
