// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"newreleases.io/newreleases"
)

func (c *command) initSlackCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "slack",
		Short: "List Slack integrations",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			elements, err := c.slackChannelsService.List(ctx)
			if err != nil {
				return err
			}

			if len(elements) == 0 {
				cmd.Println("No Slack Channels found.")
				return nil
			}

			printSlackChannelsTable(cmd, elements)

			return nil
		},
		PreRunE: c.setSlackChannelsService,
	}

	if err := addClientFlags(cmd, c.config); err != nil {
		return err
	}

	c.root.AddCommand(cmd)
	return nil
}

func (c *command) setSlackChannelsService(cmd *cobra.Command, args []string) (err error) {
	if c.slackChannelsService != nil {
		return nil
	}
	client, err := c.getClient(cmd)
	if err != nil {
		return err
	}
	c.slackChannelsService = client.SlackChannels
	return nil
}

type slackChannelsService interface {
	List(ctx context.Context) (channels []newreleases.SlackChannel, err error)
}

func printSlackChannelsTable(cmd *cobra.Command, elements []newreleases.SlackChannel) {
	table := newTable(cmd.OutOrStdout())
	table.SetHeader([]string{"ID", "Workspace", "Channel"})
	for _, e := range elements {
		table.Append([]string{e.ID, e.TeamName, e.Channel})
	}
	table.Render()
}
