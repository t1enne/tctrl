package src

import (
	"fmt"
	"testing"
	"time"
)

func TestStrToDate(t *testing.T) {
	arg := "24-01-01"
	d := StrToDate(arg)
	if d.Year() != 2024 {
		t.Errorf("Year is not 2024 %d", d.Year())
	}
	if d.Month() != 1 {
		t.Errorf("Month is not %d", int(d.Month()))
	}
	if d.Day() != 1 {
		t.Errorf("Day is not %d", d.Day())
	}
}

func testDate(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.Now().UTC().Location())
}

func TestHour(t *testing.T) {
	minutesToAdd := time.Hour*8 + time.Minute*30
	date := testDate(2025, 1, 1).Add(minutesToAdd)
	want := "2025-01-01T08:30:00.000Z"
	got := date.Format(DATE_ISO_TMPL)
	if want != got {
		t.Errorf("want %s\ngot %s\n", want, got)
	}
}

func TestCountHoursOffTable(t *testing.T) {
	var tests = []struct {
		from  time.Time
		to    time.Time
		count float64
	}{
		{testDate(2025, 1, 2), EndOfDay(testDate(2025, 1, 3)), 16.0},
		{testDate(2025, 1, 1).Add(time.Hour*8 + time.Minute*30), testDate(2025, 1, 1).Add(time.Hour*16 + time.Minute*30), 7.0},
		{testDate(2025, 1, 1).Add(time.Hour*8 + time.Minute*30), testDate(2025, 1, 1).Add(time.Hour*9 + time.Minute*30), 1.0},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s want:%f", tt.from, tt.to, tt.count)
		t.Run(testname, func(t *testing.T) {
			if CountOffHours(tt.from, tt.to) != tt.count {
				t.Errorf("want %f, got %f", tt.count, CountOffHours(tt.from, tt.to))
			}
		})
	}
}

func TestDayUtilsTable(t *testing.T) {
	var tests = []struct {
		t     time.Time
		start string
		end   string
	}{
		{testDate(2006, 1, 2), "2006-01-02T00:00:00.000Z", "2006-01-02T23:59:59.999Z"},
		{testDate(2025, 1, 4), "2025-01-04T00:00:00.000Z", "2025-01-04T23:59:59.999Z"},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%s,%s", tt.start, tt.end)
		t.Run(testname, func(t *testing.T) {
			start := StartOfDay(tt.t).Format(DATE_ISO_TMPL)
			end := EndOfDay(tt.t).Format(DATE_ISO_TMPL)
			if start != tt.start {
				t.Errorf("start mismatch, got %s, want %s", start, tt.start)
			}
			if end != tt.end {
				t.Errorf("end mismatch, got %s, want %s", end, tt.end)
			}
		})
	}
}
