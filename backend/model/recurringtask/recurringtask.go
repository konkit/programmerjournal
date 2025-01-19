package recurringtask

import (
	"fmt"
	"programmerjournal-backend/model/date"
	"strings"
	"time"
)

type Freq string

const (
	MON = "MON"
	TUE = "TUE"
	WED = "WED"
	THU = "THU"
	FRI = "FRI"
	SAT = "SAT"
	SUN = "SUN"
)

var daysOfWeek = []string{MON, TUE, WED, THU, FRI, SAT, SUN}

type FreqSettings struct {
	ByWeekDay string
}

type RecurringTask struct {
	ID              uint   `json:"id" sql:"AUTO_INCREMENT" gorm:"primarykey"`
	TaskTitle       string `json:"taskTitle"`
	TaskDescription string `json:"taskDescription"`
	FreqByWeekDay   string `json:"freqByWeekDay"`
}

func (f *FreqSettings) DayWithinDate(d date.DayDate) bool {
	dow := daysOfWeek[getDayOfWeek(d)]

	var taskDoW []string
	for _, dayOfWeek := range daysOfWeek {
		if strings.Contains(f.ByWeekDay, dayOfWeek) {
			taskDoW = append(taskDoW, dayOfWeek)
		}
	}

	for _, dd := range taskDoW {
		if dd == dow {
			return true
		}
	}

	return false
}

func getDayOfWeek(d date.DayDate) int {
	// Parse the date string into a time.Time object
	t, err := time.Parse("2006-01-02", string(d.Value))
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return 0
	}

	// Get the day of the week
	weekday := int(t.Weekday()) - 1
	if weekday == -1 {
		weekday = 6
	}

	return weekday
}
