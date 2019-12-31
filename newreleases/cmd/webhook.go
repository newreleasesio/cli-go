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

func (c *command) initWebhookCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "webhook",
		Short: "List custom Webhooks",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			webhooks, err := c.webhooksService.List(ctx)
			if err != nil {
				return err
			}

			if len(webhooks) == 0 {
				cmd.Println("No Webhooks found.")
				return nil
			}

			printWebhooksTable(cmd, webhooks)

			return nil
		},
		PreRunE: c.setWebhooksService,
	}

	if err := addClientFlags(cmd, c.config); err != nil {
		return err
	}

	c.root.AddCommand(cmd)
	return nil
}

func (c *command) setWebhooksService(cmd *cobra.Command, args []string) (err error) {
	if c.webhooksService != nil {
		return nil
	}
	client, err := c.getClient(cmd)
	if err != nil {
		return err
	}
	c.webhooksService = client.Webhooks
	return nil
}

type webhooksService interface {
	List(ctx context.Context) (webhooks []newreleases.Webhook, err error)
}

func printWebhooksTable(cmd *cobra.Command, webhooks []newreleases.Webhook) {
	table := newTable(cmd.OutOrStdout())
	table.SetHeader([]string{"ID", "Name"})
	for _, e := range webhooks {
		table.Append([]string{e.ID, e.Name})
	}
	table.Render()
}
