// Copyright (c) 2022, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"
)

func (c *command) initTagListCmd(tagCmd *cobra.Command) (err error) {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Get all tags",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			tags, err := c.tagsService.List(ctx)
			if err != nil {
				return err
			}

			if len(tags) == 0 {
				cmd.Println("No tags found.")
				return nil
			}

			printTagsTable(cmd, tags)

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := addClientConfigOptions(cmd, c.config); err != nil {
				return err
			}
			return c.setTagsService(cmd, args)
		},
	}

	tagCmd.AddCommand(cmd)
	return addClientFlags(cmd)
}
