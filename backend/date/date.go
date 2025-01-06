package date

import (
	"fmt"
	"time"
)

type Date string

func Parse(value string) Date {
	return Date(value)
}

func PlusDays(today Date, i int) Date {
	layout := "2006-01-02"

	t, err := time.Parse(layout, string(today))
	if err != nil {
		fmt.Printf("Error during minusDays(): %v\n", err)
		return ""
	}

	// Subtract one day
	yesterday := t.AddDate(0, 0, i)

	return Date(yesterday.Format(layout))
}

func MinusDays(today Date, i int) Date {
	return PlusDays(today, -i)
}

func (d Date) IsAfter(t Date) (bool, error) {
	layout := "2006-01-02"

	currentDate, err := time.Parse(layout, string(d))
	if err != nil {
		return false, fmt.Errorf("error during IsAfter(): %v\n", err)
	}

	targetDate, err := time.Parse(layout, string(t))
	if err != nil {
		return false, fmt.Errorf("error during IsAfter(): %v\n", err)
	}

	return currentDate.After(targetDate), nil
}
