/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package dayoff

import (
	"os"

	"github.com/spf13/cobra"
)

var DayOffCmd = &cobra.Command{
	Use:   "dayoff",
	Short: "Handle holiday requests",
}

func Execute() {
	err := DayOffCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	DayOffCmd.AddCommand(DayOffLogCmd)
	DayOffCmd.AddCommand(DayOffAddCmd)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// chillCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// chillCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
