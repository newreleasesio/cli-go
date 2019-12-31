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

func (c *command) initDiscordCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "discord",
		Short: "List Discord integrations",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			channel, err := c.discordChannelsService.List(ctx)
			if err != nil {
				return err
			}

			if len(channel) == 0 {
				cmd.Println("No Discord Channels found.")
				return nil
			}

			printDiscordChannelsTable(cmd, channel)

			return nil
		},
		PreRunE: c.setDiscordChannelsService,
	}

	if err := addClientFlags(cmd, c.config); err != nil {
		return err
	}

	c.root.AddCommand(cmd)
	return nil
}

func (c *command) setDiscordChannelsService(cmd *cobra.Command, args []string) (err error) {
	if c.discordChannelsService != nil {
		return nil
	}
	client, err := c.getClient(cmd)
	if err != nil {
		return err
	}
	c.discordChannelsService = client.DiscordChannels
	return nil
}

type discordChannelsService interface {
	List(ctx context.Context) (channels []newreleases.DiscordChannel, err error)
}

func printDiscordChannelsTable(cmd *cobra.Command, channels []newreleases.DiscordChannel) {
	table := newTable(cmd.OutOrStdout())
	table.SetHeader([]string{"ID", "Name"})
	for _, e := range channels {
		table.Append([]string{e.ID, e.Name})
	}
	table.Render()
}
