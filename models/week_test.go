package models

import (
	"testing"
	"time"

	"github.com/sknutsen/planner/lib"
)

func TestPopulateWeekDates_SetsWeekdayFields(t *testing.T) {
	var w Week
	// 2025-01-06 is a Monday (UTC)
	dates := []string{"2025-01-06", "2025-01-07", "2025-01-08", "2025-01-09", "2025-01-10", "2025-01-11", "2025-01-12"}
	PopulateWeekDates(&w, dates)

	wantMon := mustDate(t, "2025-01-06")
	if !w.Monday.Date.Equal(wantMon) {
		t.Fatalf("Monday: want %v got %v", wantMon, w.Monday.Date)
	}
	if !w.Sunday.Date.Equal(mustDate(t, "2025-01-12")) {
		t.Fatalf("Sunday: got %v", w.Sunday.Date)
	}
}

func TestPopulateWeekDates_SkipsInvalidStrings(t *testing.T) {
	var w Week
	PopulateWeekDates(&w, []string{"not-a-date", "2025-01-06"})
	if !w.Monday.Date.Equal(mustDate(t, "2025-01-06")) {
		t.Fatalf("Monday: got %v", w.Monday.Date)
	}
}

func TestWeek_GetWeekday_Match(t *testing.T) {
	d := mustDate(t, "2025-01-08") // Wednesday
	var w Week
	w.Wednesday.Date = d

	got, err := w.GetWeekday(d)
	if err != nil {
		t.Fatal(err)
	}
	if !got.Date.Equal(d) {
		t.Fatalf("want %v got %v", d, got.Date)
	}
}

func TestWeek_GetWeekday_Mismatch(t *testing.T) {
	var w Week
	w.Monday.Date = mustDate(t, "2025-01-06")

	other := mustDate(t, "2025-01-15")
	_, err := w.GetWeekday(other)
	if err == nil {
		t.Fatal("expected error when date is not the stored weekday")
	}
}

func mustDate(t *testing.T, s string) time.Time {
	t.Helper()
	d, err := lib.StringToDate(s)
	if err != nil {
		t.Fatal(err)
	}
	return d
}
