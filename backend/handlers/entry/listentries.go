package entry

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"programmerjournal-backend/database"
	"programmerjournal-backend/model/date"
	"programmerjournal-backend/model/entry"

	"github.com/danielgtaylor/huma/v2"
)

type ListEntriesInput struct {
	Date string `path:"date" example:"2024-05-05" doc:"Day to select the list from"`
}

type ListEntriesResponse struct {
	Pending []entry.Entry `json:"pending"`
	Done    []entry.Entry `json:"done"`
}

type ListEntriesOutput struct {
	Body   ListEntriesResponse
	Status int
}

func ListEntriesHandler(api huma.API, es *database.EntryService) {
	op := huma.Operation{
		OperationID: "ListEntries",
		Method:      http.MethodGet,
		Path:        "/api/entries/list/{date}",
		Tags:        []string{"Entry"},
	}
	huma.Register(api, op, func(ctx context.Context, input *ListEntriesInput) (*ListEntriesOutput, error) {
		resp := &ListEntriesOutput{}

		dateType := date.GetDateType(input.Date)

		switch dateType {
		case date.DateTypeDay:
			dayDate, err := date.ParseDayDate(date.DateString(input.Date))
			if err != nil {
				resp.Status = http.StatusBadRequest
				return nil, err
			}

			entries, err := es.FindEntriesByDate(dayDate.Value)
			if err != nil {
				slog.Error("Error in ListEntriesHandler", "error", err)
				return nil, huma.Error500InternalServerError(err.Error())
			}

			resp.Body = splitEntries(entries)
			resp.Status = http.StatusOK
			return resp, nil
		case date.DateTypeMonth:
			monthDate, err := date.ParseMonthDate(date.DateString(input.Date))
			if err != nil {
				resp.Status = http.StatusBadRequest
				return nil, err
			}

			entries, err := es.FindEntriesByDate(monthDate.Value)
			if err != nil {
				slog.Error("Error in ListEntriesHandler", "error", err)
				return nil, huma.Error500InternalServerError(err.Error())
			}

			resp.Body = splitEntries(entries)
			resp.Status = http.StatusOK
			return resp, nil
		case date.DateTypeWeek:
			weekDate, err := date.ParseWeekDate(date.DateString(input.Date))
			if err != nil {
				resp.Status = http.StatusBadRequest
				return nil, err
			}

			entries, err := es.FindEntriesByDate(weekDate.Value)
			if err != nil {
				slog.Error("Error in ListEntriesHandler", "error", err)
				return nil, huma.Error500InternalServerError(err.Error())
			}

			resp.Body = splitEntries(entries)
			resp.Status = http.StatusOK
			return resp, nil
		default:
			resp.Status = http.StatusBadRequest
			return nil, fmt.Errorf("unrecognized date format: %s", input.Date)
		}
	})
}

func splitEntries(entries []entry.Entry) ListEntriesResponse {
	pending := []entry.Entry{}
	done := []entry.Entry{}
	for _, e := range entries {
		if e.Status == entry.StatusTaskMigrated || e.Status == entry.StatusTaskCancelled || e.Status == entry.StatusTaskSnoozed {
			done = append(done, e)
		} else {
			pending = append(pending, e)
		}
	}
	return ListEntriesResponse{
		Pending: pending,
		Done:    done,
	}
}
