/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package dayoff

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/t1enne/tctrl/src"
)

type GenWithId struct {
	ID string `json:"_id"`
}

var DayOffAddCmd = &cobra.Command{
	Use:   "add",
	Short: "add a dayoff req",
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := DayOffCmd.Flags().GetString("config")
		config := src.GetConfig(configPath)
		fromDate, toDate := src.HandleArgs(cmd)
		isUnder24hrs := toDate.Day() != fromDate.Day()
		if !isUnder24hrs {
			notes := createNotesForm()
			from, to := createHoursForm(fromDate)
			payload := src.AddDayOffPayload{
				StartDate: from.Format(src.DATE_ISO_TMPL),
				EndDate:   to.Format(src.DATE_ISO_TMPL),
				Notes:     notes,
				Hours:     src.CalcWorkedHours(from, to),
				User:      GenWithId{ID: config.User.Id},
				Status:    "Da approvare",
			}
			src.AddDayOff(payload, config)
			return
		}
		notes := createNotesForm()
		payload := src.AddDayOffPayload{
			StartDate: src.StartOfDay(fromDate).Format(src.DATE_ISO_TMPL),
			EndDate:   src.EndOfWorkingDay(toDate).Format(src.DATE_ISO_TMPL),
			Notes:     notes,
			Hours:     src.CountOffHours(src.StartOfDay(fromDate), src.EndOfDay(toDate)),
			User:      GenWithId{ID: config.User.Id},
			Status:    "Da approvare",
		}
		src.AddDayOff(payload, config)
	},
}

func init() {}
func createHoursForm(day time.Time) (time.Time, time.Time) {
	var fromTime time.Time
	var toTime time.Time
	var hours = []string{"08.30", "09.00", "09.30", "10.00", "10.30", "11.00", "11.30", "12.00", "12.30", "13.00", "14.00", "14.30", "15.00", "15.30", "16.00", "16.30", "17.00", "17.30"}
	var hrsOpts = make([]huh.Option[time.Time], len(hours))
	for i, t := range hours {
		dt, _ := time.Parse("06-01-02:15.04", day.Format("06-01-02")+":"+t)
		atHour := time.Date(dt.Year(), dt.Month(), dt.Day(), dt.Hour()-1, dt.Minute(), 0, 0, dt.UTC().Location())
		hrsOpts[i] = huh.NewOption[time.Time](t, atHour)
	}
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[time.Time]().
				Height(8).
				Title("From").
				Options(hrsOpts...).Value(&fromTime),

			huh.NewSelect[time.Time]().
				Height(8).
				Title("To").
				Options(hrsOpts...).Value(&toTime),
		),
	)
	err := f.Run()
	if err != nil {
		fmt.Println("terminated")
		os.Exit(1)
	}

	return fromTime, toTime
}
func createNotesForm() string {
	var v string
	validatefn := func(str string) error {
		if str == "" {
			return errors.New("empty notes")
		}
		return nil
	}
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Notes").
				Validate(validatefn).
				Value(&v),
		),
	)
	err := f.Run()
	if err != nil {
		fmt.Println("terminated")
		os.Exit(1)
	}
	return v
}
