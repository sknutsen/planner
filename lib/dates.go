package lib

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func StringToDate(s string) (time.Time, error) {
	return time.Parse(time.DateOnly, s)
}

func DateToString(d time.Time) string {
	return d.Format(time.DateOnly)
}

func DatesInWeek(y int, w int) []string {
	dates := []string{}

	start := WeekStart(y, w)

	dates = append(dates, DateToString(start))

	for i := 1; i < 7; i++ {
		dates = append(dates, DateToString(start.AddDate(0, 0, i)))
	}

	return dates
}

func ISOWeekFromString(w string) (int, int) {
	println(w)
	parts := strings.Split(w, "-")

	year, _ := strconv.Atoi(parts[0])
	week, _ := strconv.Atoi(parts[1])
	return year, week
}

func WeekStart(year int, week int) time.Time {
	// Start from the middle of the year:
	t := time.Date(year, 7, 1, 0, 0, 0, 0, time.UTC)

	// Roll back to Monday:
	if wd := t.Weekday(); wd == time.Sunday {
		t = t.AddDate(0, 0, -6)
	} else {
		t = t.AddDate(0, 0, -int(wd)+1)
	}

	// Difference in weeks:
	_, w := t.ISOWeek()
	t = t.AddDate(0, 0, (week-w)*7)

	return t
}

func PrevISOWeek(w string) string {
	year, week := ISOWeekFromString(w)
	return ISOWeek(WeekStart(year, week-1))
}

func NextISOWeek(w string) string {
	year, week := ISOWeekFromString(w)
	return ISOWeek(WeekStart(year, week+1))
}

func PrevDateString(d time.Time) string {
	return DateToString(d.AddDate(0, 0, -1))
}

func NextDateString(d time.Time) string {
	return DateToString(d.AddDate(0, 0, 1))
}

func ISOWeek(d time.Time) string {
	y, w := d.ISOWeek()

	return fmt.Sprintf("%d-%d", y, w)
}
