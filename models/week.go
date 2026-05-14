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

// PopulateWeekDates sets each weekday's Date from ISO date strings (YYYY-MM-DD).
func PopulateWeekDates(w *Week, dateStrings []string) {
	for _, s := range dateStrings {
		date, err := lib.StringToDate(s)
		if err != nil {
			continue
		}
		switch date.Weekday() {
		case time.Monday:
			w.Monday.Date = date
		case time.Tuesday:
			w.Tuesday.Date = date
		case time.Wednesday:
			w.Wednesday.Date = date
		case time.Thursday:
			w.Thursday.Date = date
		case time.Friday:
			w.Friday.Date = date
		case time.Saturday:
			w.Saturday.Date = date
		case time.Sunday:
			w.Sunday.Date = date
		}
	}
}
