/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/t1enne/tctrl/cmd/dayoff"
	"github.com/t1enne/tctrl/src"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tctrl",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(dayoff.DayOffCmd)
	rootCmd.PersistentFlags().StringP("config", "c", src.GetConfigPath(), "config path")
	rootCmd.PersistentFlags().StringP("exact", "e", "", "operate on <date>")
	rootCmd.PersistentFlags().StringP("from", "f", "", "operate from <date>")
	rootCmd.PersistentFlags().StringP("to", "t", "", "operate up to <date>")
}
