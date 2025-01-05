package src

import (
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
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

func HandleArgs(cmd *cobra.Command) (time.Time, time.Time) {
	exactArg, _ := cmd.Flags().GetString("exact")
	fromArg, _ := cmd.Flags().GetString("from")
	toArg, _ := cmd.Flags().GetString("to")
	// NO ARGS
	if exactArg == "" && fromArg == "" && toArg == "" {
		n := time.Now()
		return StartOfDay(n), EndOfDay(n)
	}
	// BOTH EXACT AND FROM/TO
	if exactArg != "" && (fromArg != "" || toArg != "") {
		log.Panicln("Cannot set both --exact and --from or --to")
	}
	// ONLY EXACT
	if exactArg != "" {
		return StartOfDay(StrToDate(exactArg)), EndOfDay(StrToDate(exactArg))
	}
	// FROM
	fromDate := StrToDate(fromArg)
	var toDate time.Time
	if toArg != "" {
		toDate = StrToDate(toArg)
	} else {
		toDate = time.Now()
	}
	return fromDate, toDate
}

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

func ToFloat(s string) float64 {
	value, e := strconv.ParseFloat(s, 64)
	if e != nil {
		log.Panicln("Failed to conver value " + s)
	}
	return value
}

func CountOffHours(from time.Time, to time.Time) float64 {
	count := 0.0
	for dayIter := from; dayIter.Before(to); dayIter = StartOfWorkingDay(dayIter.Add(time.Hour * 24)) {
		if !IsWeekend(dayIter) {
			count += 8
		}
	}
	return count
}

func CalcWorkedHours(from time.Time, to time.Time) float64 {
	if from.Day() != to.Day() {
		log.Panicln("from and to must be on the same day")
	}

	// Define key time boundaries
	startOfDay := time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
	workStart := startOfDay.Add(7*time.Hour + 30*time.Minute) // 7:30
	lunchStart := startOfDay.Add(12 * time.Hour)              // 12:00
	lunchEnd := startOfDay.Add(13 * time.Hour)                // 13:00
	workEnd := startOfDay.Add(16*time.Hour + 30*time.Minute)  // 16:30

	// Upperbound variable to keep time within valid ranges
	upperBound := func(t time.Time, min time.Time, max time.Time) time.Time {
		if t.Before(min) {
			return min
		}
		if t.After(max) {
			return max
		}
		return t
	}

	// Adjust `from` and `to` within workday boundaries
	from = upperBound(from, workStart, workEnd)
	to = upperBound(to, workStart, workEnd)

	// If the range is entirely outside of work hours
	if from.After(to) {
		return 0
	}

	// Calculate hours
	totalHours := 0.0

	// Morning session: workStart to lunchStart
	if from.Before(lunchStart) {
		morningEnd := upperBound(to, workStart, lunchStart)
		totalHours += morningEnd.Sub(from).Hours()
	}

	// Afternoon session: lunchEnd to workEnd
	if to.After(lunchEnd) {
		afternoonStart := upperBound(from, lunchEnd, workEnd)
		totalHours += to.Sub(afternoonStart).Hours()
	}

	return math.Min(totalHours, 8.0)
}

func StartOfWorkingDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 7, 30, 0, 000000000, t.UTC().Location())
}

func EndOfWorkingDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 16, 30, 0, 000000000, t.UTC().Location())
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
