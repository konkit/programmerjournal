package date

import (
	"fmt"
	"time"
)

type Date string

func Parse(value string) Date {
	return Date(value)
}

func MinusDays(today Date, i int) Date {
	layout := "2006-01-02"

	t, err := time.Parse(layout, string(today))
	if err != nil {
		fmt.Printf("Error during minusDays(): %v\n", err)
		return ""
	}

	// Subtract one day
	yesterday := t.AddDate(0, 0, -i)

	return Date(yesterday.Format(layout))
}
