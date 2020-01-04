// Copyright (c) 2020, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"
	"newreleases.io/newreleases"
)

func (c *command) initProjectListCmd(projectCmd *cobra.Command) (err error) {
	var (
		optionNamePage     = "page"
		optionNameProvider = "provider"
		optionNameOrder    = "order"
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "Get tracked projects",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			flags := cmd.Flags()
			page, err := flags.GetInt(optionNamePage)
			if err != nil {
				return err
			}
			provider, err := flags.GetString(optionNameProvider)
			if err != nil {
				return err
			}
			order, err := flags.GetString(optionNameOrder)
			if err != nil {
				return err
			}

			o := newreleases.ProjectListOptions{
				Page:     page,
				Provider: provider,
			}
			if order != "" {
				o.Order = newreleases.ProjectListOrder(order)
			}

			projects, lastPage, err := c.projectsService.List(ctx, o)
			if err != nil {
				return err
			}

			if len(projects) == 0 || err == newreleases.ErrNotFound {
				if page <= 1 {
					cmd.Println("No projects found.")
					return nil
				}
				cmd.Printf("No projects found on page %v.\n", page)
				return nil
			}

			printProjectsTable(cmd, projects)

			if page < lastPage {
				cmd.Println("More projects on the next page...")
			}

			return nil
		},
		PreRunE: c.setProjectsService,
	}

	cmd.Flags().IntP(optionNamePage, "p", 1, "page number")
	cmd.Flags().String(optionNameProvider, "", "filter by provider")
	cmd.Flags().String(optionNameOrder, "", "sort projects: updated, added, name; default updated")

	if err := addClientFlags(cmd, c.config); err != nil {
		return err
	}

	projectCmd.AddCommand(cmd)
	return nil
}
