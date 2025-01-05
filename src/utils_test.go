package src

import (
	"fmt"
	"testing"
	"time"
)

func TestStrToDate(t *testing.T) {
	arg := "2024-01-01"
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
