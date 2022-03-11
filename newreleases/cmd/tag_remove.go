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

func (c *command) initTagRemoveCmd(tagCmd *cobra.Command) (err error) {
	cmd := &cobra.Command{
		Use:   "remove TAG_ID",
		Short: "Remove a tag by its ID",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			if len(args) != 1 {
				return cmd.Help()
			}

			if err := c.tagsService.Delete(ctx, args[0]); err != nil {
				if errors.Is(err, newreleases.ErrNotFound) {
					cmd.Println("Tag not found.")
					return nil
				}
				return err
			}
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
