package date

import "testing"

func TestDateString_IsAfter(t *testing.T) {
	testCases := []struct {
		name    string
		date    DateString
		other   DateString
		want    bool
		wantErr bool
	}{
		{
			name:    "date_is_after",
			date:    DateString("2024-05-01"),
			other:   DateString("2024-04-01"),
			want:    true,
			wantErr: false,
		},
		{
			name:    "date_is_before",
			date:    DateString("2024-05-01"),
			other:   DateString("2024-06-01"),
			want:    false,
			wantErr: false,
		},
		{
			name:    "date_is_the_same",
			date:    DateString("2024-05-01"),
			other:   DateString("2024-06-01"),
			want:    false,
			wantErr: false,
		},
		{
			name:    "different_formats",
			date:    DateString("2024-05-01"),
			other:   DateString("2024-06"),
			want:    false,
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := tc.date.IsAfter(tc.other)
			if tc.wantErr && err == nil || !tc.wantErr && err != nil {
				t.Error(err)
			}
			if res != tc.want {
				t.Errorf("%s.IsAfter(%s) got %v, want %v", tc.date, tc.other, res, tc.want)
			}
		})
	}

}

func TestWeekDate_PlusWeek(t *testing.T) {
	testCases := []struct {
		name     string
		date     WeekDate
		weeks    int
		expected WeekDate
	}{
		{
			name:     "plus_one_week",
			date:     WeekDate{Value: "2024-W01"},
			weeks:    1,
			expected: WeekDate{Value: "2024-W02"},
		},
		{
			name:     "plus_ten_weeks",
			date:     WeekDate{Value: "2024-W01"},
			weeks:    10,
			expected: WeekDate{Value: "2024-W11"},
		},
		{
			name:     "cross_year_boundary",
			date:     WeekDate{Value: "2024-W52"},
			weeks:    1,
			expected: WeekDate{Value: "2025-W01"},
		},
		{
			name:     "minus_one_week",
			date:     WeekDate{Value: "2024-W02"},
			weeks:    -1,
			expected: WeekDate{Value: "2024-W01"},
		},
		{
			name:     "minus_cross_year_boundary",
			date:     WeekDate{Value: "2025-W01"},
			weeks:    -1,
			expected: WeekDate{Value: "2024-W52"},
		},
		{
			name:     "leap_year_week_53",
			date:     WeekDate{Value: "2020-W53"},
			weeks:    1,
			expected: WeekDate{Value: "2021-W01"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.date.PlusWeek(tc.weeks)
			if result != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}

func TestParseWeekDate(t *testing.T) {
	testCases := []struct {
		name    string
		input   DateString
		want    WeekDate
		wantErr bool
	}{
		{
			name:    "valid_week_date",
			input:   "2024-W01",
			want:    WeekDate{Value: "2024-W01"},
			wantErr: false,
		},
		{
			name:    "valid_week_date_52",
			input:   "2024-W52",
			want:    WeekDate{Value: "2024-W52"},
			wantErr: false,
		},
		{
			name:    "invalid_format_missing_W",
			input:   "2024-01",
			want:    WeekDate{},
			wantErr: true,
		},
		{
			name:    "invalid_format_wrong_length",
			input:   "2024-W1",
			want:    WeekDate{},
			wantErr: true,
		},
		{
			name:    "invalid_format_letters",
			input:   "2024-WXX",
			want:    WeekDate{},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := ParseWeekDate(tc.input)
			if (err != nil) != tc.wantErr {
				t.Errorf("ParseWeekDate() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if got != tc.want {
				t.Errorf("ParseWeekDate() got = %v, want %v", got, tc.want)
			}
		})
	}
}
