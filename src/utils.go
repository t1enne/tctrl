package src

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Year: "2006" "06"
// Month: "Jan" "January" "01" "1"
// Day of the week: "Mon" "Monday"
// Day of the month: "2" "_2" "02"
// Day of the year: "__2" "002"
// Hour: "15" "3" "03" (PM or AM)
// Minute: "4" "04"
// Second: "5" "05"
// AM/PM mark: "PM"
var DATE_ISO_TMPL = "2006-01-02T15:04:05.000Z"
var DATE_READABLE_TMPL = "Mon, 06 Jan 02 (15:04)"

func GetConfigPath() string {
	configDir, err := os.UserConfigDir()
	if err != nil {
		configDir, _ = os.UserHomeDir()
	}
	fullPath := fmt.Sprintf("%s/%s", configDir, "tcontrol.json")
	_, err = os.Open(fullPath)
	if err != nil {
		fmt.Println("Failed to find config file. Use login cmd to create one")
		os.Exit(1)
	}
	return fullPath
}

/**
Accepts a YYYY-MM-DD string and returns a time.Time
*/
func StrToDate(dateStr string) time.Time {
	parsedDate, err := time.Parse("06-01-02", strings.TrimSpace(dateStr))
	if err != nil {
		log.Panicf("Error parsing date: %s", dateStr)
	}
	return parsedDate
}

func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 000000000, t.UTC().Location())
}

func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 999000000, t.UTC().Location())
}

func IsWeekend(t time.Time) bool {
	return t.Weekday() == time.Saturday || t.Weekday() == time.Sunday
}

func FmtDate(t time.Time) string {
	y := t.Year()
	m := int(t.Month())
	d := t.Day()
	wd := t.Weekday().String()
	isToday := time.Now().Format("2006-01-02") == t.Format("2006-01-02")
	if isToday {
		return fmt.Sprintf("%d-%02d-%02d %s *", y, m, d, wd)
	} else {
		return fmt.Sprintf("%d-%02d-%02d %s", y, m, d, wd)
	}

}

func PrintDay(t time.Time, s lipgloss.Style) {
	fmt.Println(s.Render(FmtDate(t)))
}

func FmtHours(h UserHours) string {
	// 8.00  [ATM FE 2024]  (TAOV settembre - dicembre) {svi}
	return fmt.Sprintf("%s [%s] (%s) {%s}", h.Hours, h.Release.Project.Name, h.Release.Name, strings.ToLower(h.HoursTag.Name)[:3])
}

func PrintHours(h UserHours, s lipgloss.Style) {
	fmt.Println(s.Render(strings.Repeat(" ", 11) + FmtHours(h)))
}

func FmtDayOff(dayoff DayOff) string {
	start, _ := time.Parse(DATE_ISO_TMPL, dayoff.StartDate)
	end, _ := time.Parse(DATE_ISO_TMPL, dayoff.EndDate)
	return fmt.Sprintf("%s - %s [%s] %s (%s)", start.Format(DATE_READABLE_TMPL), end.Format(DATE_READABLE_TMPL), dayoff.Hours, dayoff.Notes, dayoff.Status)
}

func PrintDayOff(doff DayOff, s lipgloss.Style) {
	fmt.Println(s.Render(FmtDayOff(doff)))
}

func Filter[T any](slice []T, predicate func(T) bool) []T {
	var result []T
	for _, v := range slice {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}
