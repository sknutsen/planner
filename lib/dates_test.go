package lib

import (
	"testing"
	"time"
)

func TestStringToDate_DateOnly(t *testing.T) {
	d, err := StringToDate("2025-06-15")
	if err != nil {
		t.Fatal(err)
	}
	if d.Year() != 2025 || d.Month() != time.June || d.Day() != 15 {
		t.Fatalf("got %v", d)
	}
}

func TestDateToString_RoundTrip(t *testing.T) {
	d, err := StringToDate("2024-02-29")
	if err != nil {
		t.Fatal(err)
	}
	if got := DateToString(d); got != "2024-02-29" {
		t.Fatalf("want 2024-02-29 got %q", got)
	}
}

func TestStripDateString(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"2025-01-01", "2025-01-01"},
		{"2025-01-01T12:34:56Z", "2025-01-01"},
		{"2025-01-01T00:00:00.000Z", "2025-01-01"},
	}
	for _, tc := range cases {
		if got := StripDateString(tc.in); got != tc.want {
			t.Fatalf("StripDateString(%q) = %q want %q", tc.in, got, tc.want)
		}
	}
}

func TestISOWeek_Format(t *testing.T) {
	// 2025-01-06 is Monday of ISO week 2 in 2025 (year may differ for Jan 1 edge cases)
	d := time.Date(2025, 1, 6, 0, 0, 0, 0, time.UTC)
	got := ISOWeek(d)
	if got != "2025-2" {
		t.Fatalf("ISOWeek: got %q", got)
	}
}

func TestDatesInWeek_SevenDays(t *testing.T) {
	dates := DatesInWeek(2025, 2)
	if len(dates) != 7 {
		t.Fatalf("want 7 dates got %d: %v", len(dates), dates)
	}
	// First day should parse as Monday of that ISO week
	first, err := StringToDate(dates[0])
	if err != nil {
		t.Fatal(err)
	}
	if first.Weekday() != time.Monday {
		t.Fatalf("first day weekday: want Monday got %v", first.Weekday())
	}
}
