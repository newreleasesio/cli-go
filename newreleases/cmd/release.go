// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
	"jaytaylor.com/html2text"
	"newreleases.io/newreleases"
)

func (c *command) initReleaseCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "release",
		Short: "Get project releases and release notes",
	}

	if err := c.initReleaseListCmd(cmd); err != nil {
		return err
	}
	if err := c.initReleaseGetCmd(cmd); err != nil {
		return err
	}
	if err := c.initReleaseGetLatestCmd(cmd); err != nil {
		return err
	}
	if err := c.initReleaseNoteCmd(cmd); err != nil {
		return err
	}

	c.root.AddCommand(cmd)
	return nil
}

func (c *command) initReleaseListCmd(releaseCmd *cobra.Command) (err error) {
	optionNamePage := "page"

	cmd := &cobra.Command{
		Use:   "list [PROVIDER PROJECT_NAME] | [PROJECT_ID]",
		Short: "Get project releases",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			page, err := cmd.Flags().GetInt(optionNamePage)
			if err != nil {
				return err
			}

			var releases []newreleases.Release
			var lastPage int
			switch len(args) {
			case 1:
				releases, lastPage, err = c.releasesService.ListByProjectID(ctx, args[0], page)
			case 2:
				releases, lastPage, err = c.releasesService.ListByProjectName(ctx, args[0], args[1], page)
			default:
				return cmd.Help()
			}
			if err != nil {
				return err
			}

			if len(releases) == 0 || err == newreleases.ErrNotFound {
				if page <= 1 {
					cmd.Println("No releases found.")
					return nil
				}
				cmd.Printf("No releases found on page %v.\n", page)
				return nil
			}

			printReleasesTable(cmd, releases)

			if page < lastPage {
				cmd.Println("More releases on the next page...")
			}

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := addClientConfigOptions(cmd, c.config); err != nil {
				return err
			}
			return c.setReleasesService(cmd, args)
		},
	}

	cmd.Flags().IntP(optionNamePage, "p", 1, "page number")

	releaseCmd.AddCommand(cmd)
	return addClientFlags(cmd)
}

func (c *command) initReleaseGetCmd(releaseCmd *cobra.Command) (err error) {
	cmd := &cobra.Command{
		Use:   "get [provider project_name] | [project_id] version",
		Short: "Get a specific project release",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			var release *newreleases.Release
			switch len(args) {
			case 2:
				release, err = c.releasesService.GetByProjectID(ctx, args[0], args[1])
			case 3:
				release, err = c.releasesService.GetByProjectName(ctx, args[0], args[1], args[2])
			default:
				return cmd.Help()
			}
			if err != nil {
				return err
			}

			if release == nil {
				cmd.Println("Release not found.")
				return nil
			}

			printRelease(cmd, release)
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := addClientConfigOptions(cmd, c.config); err != nil {
				return err
			}
			return c.setReleasesService(cmd, args)
		},
	}

	releaseCmd.AddCommand(cmd)
	return addClientFlags(cmd)
}

func (c *command) initReleaseGetLatestCmd(releaseCmd *cobra.Command) (err error) {
	cmd := &cobra.Command{
		Use:   "get-latest [provider project_name] | [project_id]",
		Short: "Get the latest non-excluded project release",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			var release *newreleases.Release
			switch len(args) {
			case 1:
				release, err = c.releasesService.GetLatestByProjectID(ctx, args[0])
			case 2:
				release, err = c.releasesService.GetLatestByProjectName(ctx, args[0], args[1])
			default:
				return cmd.Help()
			}
			if err != nil {
				return err
			}

			if release == nil {
				cmd.Println("Release not found.")
				return nil
			}

			printRelease(cmd, release)
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := addClientConfigOptions(cmd, c.config); err != nil {
				return err
			}
			return c.setReleasesService(cmd, args)
		},
	}

	releaseCmd.AddCommand(cmd)
	return addClientFlags(cmd)
}

func (c *command) initReleaseNoteCmd(releaseCmd *cobra.Command) (err error) {
	cmd := &cobra.Command{
		Use:   "note [PROVIDER PROJECT_NAME] | [PROJECT_ID] version",
		Short: "Get a project release note",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			var releaseNote *newreleases.ReleaseNote
			switch len(args) {
			case 2:
				releaseNote, err = c.releasesService.GetNoteByProjectID(ctx, args[0], args[1])
			case 3:
				releaseNote, err = c.releasesService.GetNoteByProjectName(ctx, args[0], args[1], args[2])
			default:
				return cmd.Help()
			}
			if err != nil {
				return err
			}

			if releaseNote == nil {
				cmd.Println("Release note not found.")
				return nil
			}

			printReleaseNote(cmd, releaseNote)
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := addClientConfigOptions(cmd, c.config); err != nil {
				return err
			}
			return c.setReleasesService(cmd, args)
		},
	}

	releaseCmd.AddCommand(cmd)
	return addClientFlags(cmd)
}

func (c *command) setReleasesService(cmd *cobra.Command, args []string) (err error) {
	if c.releasesService != nil {
		return nil
	}
	client, err := c.getClient(cmd)
	if err != nil {
		return err
	}
	c.releasesService = client.Releases
	return nil
}

type releasesService interface {
	ListByProjectID(ctx context.Context, projectID string, page int) (releases []newreleases.Release, lastPage int, err error)
	ListByProjectName(ctx context.Context, provider, projectName string, page int) (releases []newreleases.Release, lastPage int, err error)
	GetByProjectID(ctx context.Context, projectID, version string) (release *newreleases.Release, err error)
	GetByProjectName(ctx context.Context, provider, projectName, version string) (release *newreleases.Release, err error)
	GetLatestByProjectID(ctx context.Context, projectID string) (release *newreleases.Release, err error)
	GetLatestByProjectName(ctx context.Context, provider, projectName string) (release *newreleases.Release, err error)
	GetNoteByProjectID(ctx context.Context, projectID string, version string) (release *newreleases.ReleaseNote, err error)
	GetNoteByProjectName(ctx context.Context, provider, projectName string, version string) (release *newreleases.ReleaseNote, err error)
}

func printReleasesTable(cmd *cobra.Command, releases []newreleases.Release) {
	table := newTable(cmd.OutOrStdout())
	table.SetHeader([]string{"Version", "Date", "Pre-Release", "Has Note", "Updated", "Excluded", "CVE"})
	for _, r := range releases {
		table.Append([]string{r.Version, r.Date.Local().String(), yesNo(r.IsPrerelease), yesNo(r.HasNote), yesNo(r.IsUpdated), yesNo(r.IsExcluded), yesNo(len(r.CVE) > 0)})
	}
	table.Render()
}

func printRelease(cmd *cobra.Command, r *newreleases.Release) {
	table := newTable(cmd.OutOrStdout())
	table.Append([]string{"Version:", r.Version})
	table.Append([]string{"Date:", r.Date.Local().String()})

	if r.IsPrerelease {
		table.Append([]string{"Pre-Release:", "yes"})
	}
	if r.HasNote {
		table.Append([]string{"Has Note:", "yes"})
	}
	if r.IsUpdated {
		table.Append([]string{"Updated:", "yes"})
	}
	if r.IsExcluded {
		table.Append([]string{"Excluded:", "yes"})
	}
	if len(r.CVE) > 0 {
		table.Append([]string{"CVE:", strings.Join(r.CVE, ", ")})
	}
	table.Render()
}

func printReleaseNote(cmd *cobra.Command, n *newreleases.ReleaseNote) {
	if n.Title != "" {
		cmd.Println(strings.TrimSpace(n.Title))
		cmd.Println()
	}
	if n.Message != "" {
		message, err := html2text.FromString(n.Message, html2text.Options{PrettyTables: true})
		if err != nil {
			panic(err)
		}
		cmd.Println(strings.TrimSpace(message))
		cmd.Println()
	}
	if n.URL != "" {
		cmd.Println(strings.TrimSpace(n.URL))
		cmd.Println()
	}
}
