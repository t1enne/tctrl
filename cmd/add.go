/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"sort"
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
		fromDate, _ := src.HandleArgs(cmd)
		fromIso := src.StartOfDay(fromDate).Format(src.DATE_ISO_TMPL)

		tags := src.GetActiveTags(config)
		tag, _ := createTagsForm(tags)

		projs := src.GetProjects(config)
		releases := src.GetReleases(config)
		release, _ := createReleasesForm(projs, releases)

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

func createReleasesForm(projs []src.Project, releases []src.Release) (src.Release, error) {
	projectMap := make(map[string]src.Project)
	for _, proj := range projs {
		projectMap[proj.ID] = proj
	}

	sort.Slice(releases, func(i, j int) bool {
		projI := projectMap[releases[i].ProjectID]
		projJ := projectMap[releases[j].ProjectID]
		return projI.Name < projJ.Name
	})

	var v src.Release
	releaseOptions := make([]huh.Option[src.Release], len(releases))
	for i, release := range releases {
		proj := projectMap[release.ProjectID]
		releaseOptions[i] = huh.NewOption[src.Release](fmt.Sprintf("[%s] - %s", proj.Name, release.Name), release)
	}

	f := huh.NewForm(
		huh.NewGroup(
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
