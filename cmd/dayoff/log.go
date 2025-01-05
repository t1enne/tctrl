/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package dayoff

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/t1enne/tctrl/src"
)

// rootCmd represents the base command when called without any subcommands
var DayOffLogCmd = &cobra.Command{
	Use:   "log",
	Short: "Log holiday requests",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := DayOffCmd.Flags().GetString("config")
		config := src.GetConfig(configPath)
		// // fromDate, toDate := handleArgs(cmd)
		// // fromIso := src.StartOfDay(fromDate).Format(src.DATE_ISO_TMPL)
		// // toIso := src.EndOfDay(toDate).Format(src.DATE_ISO_TMPL)
		p := fmt.Sprintf(`{"relations":["user","hoursTag"],"fullSearchCols":["notes"],"where":[{"userId":"%s"}],"pagination":false,"order":{"startDate":"ASC"}}`, config.User.Id)
		daysoff := src.GetDayOff(p, config)
		for _, dayoff := range daysoff {
			src.PrintDayOff(dayoff, src.WorkedStyle)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.

func init() {
}
