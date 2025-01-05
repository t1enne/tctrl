/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/t1enne/tctrl/src"
)

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "log worked hours",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
	`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := rootCmd.Flags().GetString("config")
		config := src.GetConfig(configPath)
		fromDate, toDate := src.HandleArgs(cmd)
		fromIso := src.StartOfDay(fromDate).Format(src.DATE_ISO_TMPL)
		toIso := src.EndOfDay(toDate).Format(src.DATE_ISO_TMPL)
		p := fmt.Sprintf(`{ 
			"relations": ["release", "release.project", "release.project.customer", "hoursTag"], 
			"where": {"userId": "%s", "date": {"_fn": 17, "args": ["%s", "%s"]}}
			}`, config.User.Id, fromIso, toIso)
		hours := src.GetWorkedHours(p, config)
		dateToHours := make(map[string][]src.UserHours)
		for _, h := range hours {
			_, ok := dateToHours[h.Date[:10]]
			if ok == false {
				dateToHours[h.Date[:10]] = make([]src.UserHours, 0)
			}
			dateToHours[h.Date[:10]] = append(dateToHours[h.Date[:10]], h)
		}
		for dayIter := fromDate; dayIter.Before(toDate); dayIter = dayIter.Add(time.Hour * 24) {
			worked, ok := dateToHours[dayIter.Format("2006-01-02")]
			if ok || src.IsWeekend(dayIter) {
				src.PrintDay(dayIter, src.WorkedStyle)
			} else {
				src.PrintDay(dayIter, src.EmptyStyle)
			}
			if ok {
				for _, w := range worked {
					src.PrintHours(w, src.WorkedStyle)
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(logCmd)
}
