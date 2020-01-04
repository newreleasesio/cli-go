// Copyright (c) 2020, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"newreleases.io/newreleases"
)

func (c *command) initProjectSearchCmd(projectCmd *cobra.Command) (err error) {
	var optionNameProvider = "provider"

	cmd := &cobra.Command{
		Use:   "search NAME",
		Short: "Search tracked projects by name",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			flags := cmd.Flags()
			provider, err := flags.GetString(optionNameProvider)
			if err != nil {
				return err
			}

			if len(args) == 0 {
				return cmd.Help()
			}

			projects, err := c.projectsService.Search(ctx, strings.Join(args, " "), provider)
			if err != nil {
				return err
			}

			if len(projects) == 0 || err == newreleases.ErrNotFound {
				cmd.Println("No projects found.")
				return nil
			}

			printProjectsTable(cmd, projects)
			return nil
		},
		PreRunE: c.setProjectsService,
	}

	cmd.Flags().String(optionNameProvider, "", "filter by provider")

	if err := addClientFlags(cmd, c.config); err != nil {
		return err
	}

	projectCmd.AddCommand(cmd)
	return nil
}
