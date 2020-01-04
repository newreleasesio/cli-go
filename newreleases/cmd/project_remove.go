// Copyright (c) 2020, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"
	"newreleases.io/newreleases"
)

func (c *command) initProjectRemoveCmd(projectCmd *cobra.Command) (err error) {
	cmd := &cobra.Command{
		Use:   "remove [PROVIDER PROJECT_NAME] | [PROJECT_ID]",
		Short: "Remove a tracked project",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			switch len(args) {
			case 1:
				err = c.projectsService.DeleteByID(ctx, args[0])
			case 2:
				err = c.projectsService.DeleteByName(ctx, args[0], args[1])
			default:
				return cmd.Help()
			}

			if err == newreleases.ErrNotFound {
				cmd.Println("Project not found.")
				return nil
			}

			return err
		},
		PreRunE: c.setProjectsService,
	}

	if err := addClientFlags(cmd, c.config); err != nil {
		return err
	}

	projectCmd.AddCommand(cmd)
	return nil
}
