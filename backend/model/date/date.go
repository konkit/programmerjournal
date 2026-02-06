package date

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type DateString string

func (d DateString) IsAfter(other DateString) (bool, error) {
	if len(d) != len(other) {
		return false, fmt.Errorf("date strings are not in the same format")
	}

	return strings.Compare(string(d), string(other)) > 0, nil
}

type DayDate struct {
	Value DateString
}

func ParseDayDate(value DateString) (DayDate, error) {
	match, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, string(value))
	if match {
		return DayDate{value}, nil
	} else {
		return DayDate{}, fmt.Errorf("unrecognized date format: %v", value)
	}
}

func (d DayDate) IsAfter(t DayDate) (bool, error) {
	layout := "2006-01-02"

	currentDate, err := time.Parse(layout, string(d.Value))
	if err != nil {
		return false, fmt.Errorf("error during IsAfter(): %v\n", err)
	}

	targetDate, err := time.Parse(layout, string(t.Value))
	if err != nil {
		return false, fmt.Errorf("error during IsAfter(): %v\n", err)
	}

	return currentDate.After(targetDate), nil
}

func (d DayDate) PlusDays(i int) DayDate {
	layout := "2006-01-02"

	t, err := time.Parse(layout, string(d.Value))
	if err != nil {
		fmt.Printf("Error during minusDays(): %v\n", err)
		return DayDate{}
	}

	// Subtract one day
	yesterday := t.AddDate(0, 0, i)
	return DayDate{DateString(yesterday.Format(layout))}
}

func (d DayDate) MinusDays(i int) DayDate {
	return d.PlusDays(-i)
}

type MonthDate struct {
	Value DateString
}

func ParseMonthDate(value DateString) (MonthDate, error) {
	match, _ := regexp.MatchString(`^\d{4}-\d{2}$`, string(value))
	if match {
		return MonthDate{value}, nil
	} else {
		return MonthDate{}, fmt.Errorf("unrecognized date format: %v", value)
	}
}

func (d MonthDate) IsAfter(t MonthDate) (bool, error) {
	layout := "2006-01-02"

	currentDate, err := time.Parse(layout, string(d.Value+"-01"))
	if err != nil {
		return false, fmt.Errorf("error during IsAfter(): %v\n", err)
	}

	targetDate, err := time.Parse(layout, string(t.Value+"-01"))
	if err != nil {
		return false, fmt.Errorf("error during IsAfter(): %v\n", err)
	}

	return currentDate.After(targetDate), nil
}

func (d MonthDate) PlusMonth(i int) MonthDate {
	layout := "2006-01-02"

	dayDate := string(d.Value) + "-01"
	t, err := time.Parse(layout, dayDate)
	if err != nil {
		fmt.Printf("Error during minusDays(): %v\n", err)
		return MonthDate{}
	}

	resultMonth := t.AddDate(0, i, 0)

	return MonthDate{DateString(resultMonth.Format(layout)[:7])}
}

func (d MonthDate) MinusMonth(i int) MonthDate {
	return d.PlusMonth(-i)
}

type WeekDate struct {
	Value DateString
}

func ParseWeekDate(value DateString) (WeekDate, error) {
	match, _ := regexp.MatchString(`^\d{4}-W\d{2}$`, string(value))
	if match {
		return WeekDate{value}, nil
	} else {
		return WeekDate{}, fmt.Errorf("unrecognized date format: %v", value)
	}
}

type DateType string

const (
	DateTypeDay          DateType = "day"
	DateTypeMonth        DateType = "month"
	DateTypeWeek         DateType = "week"
	DateTypeUnrecognized DateType = "unrecognized"
)

func GetDateType(str string) DateType {
	if isDayDate, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, str); isDayDate {
		return DateTypeDay
	}

	if isMonthDate, _ := regexp.MatchString(`^\d{4}-\d{2}$`, str); isMonthDate {
		return DateTypeMonth
	}

	if isWeekDate, _ := regexp.MatchString(`^\d{4}-W\d{2}$`, str); isWeekDate {
		return DateTypeWeek
	}

	return DateTypeUnrecognized
}
