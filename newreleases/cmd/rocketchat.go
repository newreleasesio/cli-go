// Copyright (c) 2021, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"newreleases.io/newreleases"
)

func (c *command) initRocketchatCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "rocketchat",
		Short: "List Rocket.Chat integrations",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			webhooks, err := c.rocketchatWebhooksService.List(ctx)
			if err != nil {
				return err
			}

			if len(webhooks) == 0 {
				cmd.Println("No Rocket.Chat Webhooks found.")
				return nil
			}

			printWebhooksTable(cmd, webhooks)

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := addClientConfigOptions(cmd, c.config); err != nil {
				return err
			}
			return c.setRocketchatWebhooksService(cmd, args)
		},
	}

	c.root.AddCommand(cmd)
	return addClientFlags(cmd)
}

func (c *command) setRocketchatWebhooksService(cmd *cobra.Command, args []string) (err error) {
	if c.rocketchatWebhooksService != nil {
		return nil
	}
	client, err := c.getClient(cmd)
	if err != nil {
		return err
	}
	c.rocketchatWebhooksService = client.RocketchatWebhooks
	return nil
}

type rocketchatWebhooksService interface {
	List(ctx context.Context) (webhooks []newreleases.Webhook, err error)
}
