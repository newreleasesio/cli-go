// Copyright (c) 2022, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"
)

func (c *command) initTagAddCmd(tagCmd *cobra.Command) (err error) {
	cmd := &cobra.Command{
		Use:   "add TAG_NAME",
		Short: "Add a new tag with a provided name",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			if len(args) != 1 {
				return cmd.Help()
			}

			tag, err := c.tagsService.Add(ctx, args[0])
			if err != nil {
				return err
			}

			printTag(cmd, tag)
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
