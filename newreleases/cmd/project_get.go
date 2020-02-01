// Copyright (c) 2020, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"
	"newreleases.io/newreleases"
)

func (c *command) initProjectGetCmd(projectCmd *cobra.Command) (err error) {
	cmd := &cobra.Command{
		Use:   "get [PROVIDER PROJECT_NAME] | [PROJECT_ID]",
		Short: "Get information about a tracked project",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			var project *newreleases.Project
			switch len(args) {
			case 1:
				project, err = c.projectsService.GetByID(ctx, args[0])
			case 2:
				project, err = c.projectsService.GetByName(ctx, args[0], args[1])
			default:
				return cmd.Help()
			}
			if err != nil {
				return err
			}

			if project == nil || err == newreleases.ErrNotFound {
				cmd.Println("Project not found.")
				return nil
			}

			printProject(cmd, project)
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := addClientFlags(cmd, c.config); err != nil {
				return err
			}
			return c.setProjectsService(cmd, args)
		},
	}

	projectCmd.AddCommand(cmd)
	return nil
}
