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
