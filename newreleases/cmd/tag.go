// Copyright (c) 2022, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"newreleases.io/newreleases"
)

func (c *command) initTagCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Manage tags",
	}

	if err := c.initTagListCmd(cmd); err != nil {
		return err
	}
	if err := c.initTagGetCmd(cmd); err != nil {
		return err
	}
	if err := c.initTagAddCmd(cmd); err != nil {
		return err
	}
	if err := c.initTagUpdateCmd(cmd); err != nil {
		return err
	}
	if err := c.initTagRemoveCmd(cmd); err != nil {
		return err
	}

	c.root.AddCommand(cmd)
	return nil
}

func (c *command) setTagsService(cmd *cobra.Command, args []string) error {
	if c.tagsService != nil {
		return nil
	}
	client, err := c.getClient(cmd)
	if err != nil {
		return err
	}
	c.tagsService = client.Tags
	return nil
}

type tagsService interface {
	List(ctx context.Context) ([]newreleases.Tag, error)
	Get(ctx context.Context, id string) (*newreleases.Tag, error)
	Add(ctx context.Context, name string) (*newreleases.Tag, error)
	Update(ctx context.Context, id, name string) (*newreleases.Tag, error)
	Delete(ctx context.Context, id string) error
}

func printTagsTable(cmd *cobra.Command, tags []newreleases.Tag) {
	table := newTable(cmd.OutOrStdout())
	table.SetHeader([]string{"ID", "Name"})
	for _, tag := range tags {
		table.Append([]string{tag.ID, tag.Name})
	}
	table.Render()
}

func printTag(cmd *cobra.Command, t *newreleases.Tag) {
	table := newTable(cmd.OutOrStdout())
	table.Append([]string{"ID:", t.ID})
	table.Append([]string{"Name:", t.Name})
	table.Render()
}
