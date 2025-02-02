package entry

import (
	"context"
	"fmt"
	"github.com/danielgtaylor/huma/v2"
	"gorm.io/gorm"
	"net/http"
	"programmerjournal-backend/database"
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

func ListEntriesHandler(api huma.API, db *gorm.DB) {
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

			entries, err := ListDayEntries(db, dayDate)
			if err != nil {
				return nil, err
			}

			resp.Body = entries
			resp.Status = http.StatusOK
			return resp, nil
		case date.DateTypeMonth:
			monthDate, err := date.ParseMonthDate(date.DateString(input.Date))
			if err != nil {
				resp.Status = http.StatusBadRequest
				return nil, err
			}

			entries, err := ListMonthEntries(db, monthDate)
			if err != nil {
				return nil, err
			}

			resp.Body = entries
			resp.Status = http.StatusOK
			return resp, nil
		default:
			resp.Status = http.StatusBadRequest
			return nil, fmt.Errorf("unrecognized date format: %s", input.Date)
		}
	})
}

func ListDayEntries(db *gorm.DB, date date.DayDate) ([]entry.Entry, error) {
	return database.FindEntriesByDate(db, date.Value)
}

func ListMonthEntries(db *gorm.DB, date date.MonthDate) ([]entry.Entry, error) {
	return database.FindEntriesByDate(db, date.Value)
}
