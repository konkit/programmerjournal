package recurringtask

import (
	"programmerjournal-backend/model/date"
	"testing"
)

func TestDayWithinDate(t *testing.T) {
	testCases := []struct {
		name      string
		date      date.DayDate
		byWeekDay string
		want      bool
	}{
		{
			name:      "basic_scenario",
			date:      date.DayDate{"2025-01-13"},
			byWeekDay: "MON,TUE",
			want:      true,
		},
		{
			name:      "out_of_range",
			date:      date.DayDate{"2025-01-13"},
			byWeekDay: "THU,FRI",
			want:      false,
		},
		{
			name:      "sunday_date_monday_byweekday",
			date:      date.DayDate{"2025-01-12"},
			byWeekDay: "MON",
			want:      false,
		},
		{
			name:      "monday_date_sunday_byweekday",
			date:      date.DayDate{"2025-01-13"},
			byWeekDay: "SUN",
			want:      false,
		},
		{
			name:      "monday_date_monday_byweekday",
			date:      date.DayDate{"2025-01-13"},
			byWeekDay: "MON",
			want:      true,
		},
		{
			name:      "sunday_date_sunday_byweekday",
			date:      date.DayDate{"2025-01-12"},
			byWeekDay: "SUN",
			want:      true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			freq := FreqSettings{ByWeekDay: tc.byWeekDay}
			res := freq.DayWithinDate(tc.date)

			if res != tc.want {
				t.Errorf("got %t, want %t", res, tc.want)
			}
		})
	}

}
