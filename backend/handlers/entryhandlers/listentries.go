package entryhandlers

import (
	"context"
	"github.com/danielgtaylor/huma/v2"
	"net/http"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"
)

type ListEntriesOutput struct {
	Body []entry.Entry
}

func ListEntries(api huma.API, taskRepo *entry.DBRepository) {
	op := huma.Operation{
		OperationID: "ListEntries",
		Method:      http.MethodGet,
		Path:        "/api/entries/list/{date}",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *struct {
		Date string `path:"date" example:"2024-05-05" doc:"Day to select the list from"`
	}) (*ListEntriesOutput, error) {
		resp := &ListEntriesOutput{}
		tasks, err := taskRepo.ListTasks(date.Parse(input.Date))
		if err != nil {
			return nil, err
		}
		resp.Body = tasks
		return resp, nil
	})
}
