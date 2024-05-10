package models

import (
	"errors"
	"time"

	"github.com/sknutsen/planner/lib"
)

type Week struct {
	ISOWeek   string
	Monday    Day
	Tuesday   Day
	Wednesday Day
	Thursday  Day
	Friday    Day
	Saturday  Day
	Sunday    Day
}

func (w *Week) GetWeekday(d time.Time) (Day, error) {
	var day Day

	switch d.Weekday() {
	case time.Monday:
		day = w.Monday
	case time.Tuesday:
		day = w.Tuesday
	case time.Wednesday:
		day = w.Wednesday
	case time.Thursday:
		day = w.Thursday
	case time.Friday:
		day = w.Friday
	case time.Saturday:
		day = w.Saturday
	case time.Sunday:
		day = w.Sunday
	}

	if day.Date != d {
		return day, errors.New("date not in current week")
	}

	return day, nil
}

func (w *Week) Prev() string {
	return lib.PrevISOWeek(w.ISOWeek)
}

func (w *Week) Next() string {
	return lib.NextISOWeek(w.ISOWeek)
}
