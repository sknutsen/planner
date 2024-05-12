package models

import (
	"fmt"
	"time"

	"github.com/sknutsen/planner/database"
	"github.com/sknutsen/planner/lib"
)

type Day struct {
	Date  time.Time
	Tasks []Task
}

func (d *Day) Prev() string {
	return lib.PrevDateString(d.Date)
}

func (d *Day) Next() string {
	return lib.NextDateString(d.Date)
}

func (d *Day) Week() string {
	return lib.ISOWeek(d.Date)
}

func (d *Day) String() string {
	year := d.Date.Year()
	month := d.Date.Month()
	day := d.Date.Day()

	var dayString string
	var monthString string

	if day < 10 {
		dayString = fmt.Sprintf("0%d", day)
	} else {
		dayString = fmt.Sprint(day)
	}

	if month < 10 {
		monthString = fmt.Sprintf("0%d", month)
	} else {
		monthString = fmt.Sprintf("%d", month)
	}

	dateAsString := fmt.Sprintf("%s.%s.%d", dayString, monthString, year)

	return dateAsString
}

func (d *Day) StringShort() string {
	month := d.Date.Month()
	day := d.Date.Day()

	var dayString string
	var monthString string

	if day < 10 {
		dayString = fmt.Sprintf("0%d", day)
	} else {
		dayString = fmt.Sprint(day)
	}

	if month < 10 {
		monthString = fmt.Sprintf("0%d", month)
	} else {
		monthString = fmt.Sprintf("%d", month)
	}

	dateAsString := fmt.Sprintf("%s.%s", dayString, monthString)

	return dateAsString
}

type DayTasksResponse struct {
	Date  string
	Tasks []database.Task
}
