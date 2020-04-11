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

func (c *command) initMicrosoftTeamsCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "microsoft-teams",
		Short: "List Microsoft Teams integrations",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			webhooks, err := c.microsoftTeamsWebhooksService.List(ctx)
			if err != nil {
				return err
			}

			if len(webhooks) == 0 {
				cmd.Println("No Microsoft Teams Webhooks found.")
				return nil
			}

			printWebhooksTable(cmd, webhooks)

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := addClientConfigOptions(cmd, c.config); err != nil {
				return err
			}
			return c.setMicrosoftTeamsWebhooksService(cmd, args)
		},
	}

	c.root.AddCommand(cmd)
	return addClientFlags(cmd)
}

func (c *command) setMicrosoftTeamsWebhooksService(cmd *cobra.Command, args []string) (err error) {
	if c.microsoftTeamsWebhooksService != nil {
		return nil
	}
	client, err := c.getClient(cmd)
	if err != nil {
		return err
	}
	c.microsoftTeamsWebhooksService = client.MicrosoftTeamsWebhooks
	return nil
}

type microsoftTeamsWebhooksService interface {
	List(ctx context.Context) (webhooks []newreleases.Webhook, err error)
}
