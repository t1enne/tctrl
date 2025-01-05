/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"log"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/spf13/cobra"
	"github.com/t1enne/tctrl/src"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Fill worked hours for dates",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		configPath, _ := rootCmd.Flags().GetString("config")
		config := src.GetConfig(configPath)
		fromDate, _ := handleArgs(cmd)
		fromIso := src.StartOfDay(fromDate).Format(src.DATE_ISO_TMPL)
		// toIso := src.EndOfDay(toDate).Format(src.DATE_ISO_TMPL)
		tags := src.GetActiveTags(config)
		tag, _ := createTagsForm(tags)

		customers := src.GetCustomers(config)
		customer, _ := createCustomersForm(customers)

		projs := src.GetProjects(customer.ID, config)
		proj, _ := createProjectsForm(projs)

		releases := src.GetReleases(proj.ID, config)
		release, _ := createReleasesForm(releases)

		notes, _ := createNotesForm()
		hours, _ := createHoursForm()

		toUpload := src.AddHoursPayload{
			Hours:      hours,
			Notes:      notes,
			Date:       fromIso,
			ReleaseId:  release.ID,
			HoursTagId: tag.ID,
			UserId:     config.User.Id,
		}

		err := spinner.New().Title("Uploading entries ...").
			Action(func() { src.AddHours(toUpload, config) }).
			Run()
		if err != nil {
			log.Panicf("failed to start spinner: %s", err)
		}
	},
}

func init() {}

func createNotesForm() (string, error) {
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
	return v, err
}

func createHoursForm() (string, error) {
	var v string
	validatefn := func(str string) error {
		if _, err := strconv.ParseFloat(v, 32); err != nil {
			return errors.New("invalid hours (accepted values are 8, 6.5) ")
		}
		return nil
	}
	f := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Hours").
				Validate(validatefn).
				Value(&v),
		),
	)
	err := f.Run()
	return v, err
}

func createReleasesForm(releases []src.Release) (src.Release, error) {
	var v src.Release
	releaseOptions := make([]huh.Option[src.Release], len(releases))
	for i, release := range releases {
		releaseOptions[i] = huh.NewOption[src.Release](release.Name, release)
	}
	f := huh.NewForm(
		huh.NewGroup(
			// Ask the user for a base burger and toppings.
			huh.NewSelect[src.Release]().
				Height(8).
				Title("Release").
				Options(
					releaseOptions...,
				).
				Value(&v),
		),
	)
	err := f.Run()
	return v, err
}

func createProjectsForm(projs []src.Project) (src.Project, error) {
	var v src.Project
	projectOptions := make([]huh.Option[src.Project], len(projs))
	for i, project := range projs {
		projectOptions[i] = huh.NewOption[src.Project](project.Name, project)
	}
	f := huh.NewForm(
		huh.NewGroup(
			// Ask the user for a base burger and toppings.
			huh.NewSelect[src.Project]().
				Height(8).
				Title("Project").
				Options(
					projectOptions...,
				).
				Value(&v),
		),
	)
	err := f.Run()
	return v, err
}

func createTagsForm(tags []src.HoursTag) (src.HoursTag, error) {
	var v src.HoursTag
	tagOptions := make([]huh.Option[src.HoursTag], len(tags))
	for i, tag := range tags {
		tagOptions[i] = huh.NewOption[src.HoursTag](tag.Name, tag)
	}
	f := huh.NewForm(
		huh.NewGroup(
			// Ask the user for a base burger and toppings.
			huh.NewSelect[src.HoursTag]().
				Height(8).
				Title("Tag").
				Options(
					tagOptions...,
				).
				Value(&v),
		),
	)
	err := f.Run()
	return v, err
}
func createCustomersForm(customers []src.Customer) (src.Customer, error) {
	var c src.Customer
	customerOptions := make([]huh.Option[src.Customer], len(customers))
	for i, c := range customers {
		customerOptions[i] = huh.NewOption[src.Customer](c.Name, c)
	}

	f := huh.NewForm(
		huh.NewGroup(
			// Ask the user for a base burger and toppings.
			huh.NewSelect[src.Customer]().
				Height(8).
				Title("Customer").
				Options(
					customerOptions...,
				).
				Value(&c),
		),
	)
	err := f.Run()
	return c, err
}
