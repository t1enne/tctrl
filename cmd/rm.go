/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
	"github.com/t1enne/tctrl/src"
)

// rmCmd represents the rm command
var rmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Delete entry",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command.`,
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
		picked, _ := createRMForm(hours)
		action := func() { deleteEntries(picked, config) }
		err := spinner.New().Title("Deleting entries ...").
			Action(action).
			Run()
		if err != nil {
			log.Panicf("failed to start spinner: %s", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}

func deleteEntries(entries []src.UserHours, c src.UserConfig) {
	for _, p := range entries {
		src.Delete("userhours/"+p.ID, "", c, nil)
	}
}

func createRMForm(hours []src.UserHours) ([]src.UserHours, error) {
	var picked []src.UserHours
	hourOptions := make([]huh.Option[src.UserHours], len(hours))
	for i, h := range hours {
		d, err := time.Parse(src.DATE_ISO_TMPL, h.Date)
		if err != nil {
			log.Panicf("failed to parse date %s\n", h.Date)
		}
		optionTxt := fmt.Sprintf("%d - %s / %s", i+1, src.FmtDate(d), src.FmtHours(h))
		hourOptions[i] = huh.NewOption[src.UserHours](optionTxt, h)
	}

	f := huh.NewForm(
		huh.NewGroup(
			// Ask the user for a base burger and toppings.
			huh.NewMultiSelect[src.UserHours]().
				Title("Pick entry to rm").
				Options(
					hourOptions...,
				).
				Value(&picked),
		),
	)
	err := f.Run()
	return picked, err
}
