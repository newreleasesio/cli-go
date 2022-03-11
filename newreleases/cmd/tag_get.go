// Copyright (c) 2022, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"newreleases.io/newreleases"
)

func (c *command) initTagGetCmd(tagCmd *cobra.Command) (err error) {
	cmd := &cobra.Command{
		Use:   "get TAG_ID",
		Short: "Get information about a tag by its ID",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			if len(args) != 1 {
				return cmd.Help()
			}

			tag, err := c.tagsService.Get(ctx, args[0])
			if err != nil {
				if errors.Is(err, newreleases.ErrNotFound) {
					cmd.Println("Tag not found.")
					return nil
				}
				return err
			}

			if tag == nil {
				cmd.Println("Tag not found.")
				return nil
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
