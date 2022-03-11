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

func (c *command) initTagUpdateCmd(tagCmd *cobra.Command) (err error) {
	var (
		optionNameName = "name"
	)

	cmd := &cobra.Command{
		Use:   "update TAG_ID",
		Short: "Update a tag referenced by its ID",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			if len(args) != 1 {
				return cmd.Help()
			}

			flags := cmd.Flags()
			name, err := flags.GetString(optionNameName)
			if err != nil {
				return err
			}

			if name == "" {
				cmd.Println("Option --name is required.")
				return nil
			}

			tag, err := c.tagsService.Update(ctx, args[0], name)
			if err != nil {
				if errors.Is(err, newreleases.ErrNotFound) {
					cmd.Println("Tag not found.")
					return nil
				}
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

	cmd.Flags().String(optionNameName, "", "tag name")

	tagCmd.AddCommand(cmd)
	return addClientFlags(cmd)
}
