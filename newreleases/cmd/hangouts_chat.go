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

func (c *command) initHangoutsChatCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "hangouts-chat",
		Short: "List Hangouts Chat integrations",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			webhooks, err := c.hangoutsChatWebhooksService.List(ctx)
			if err != nil {
				return err
			}

			if len(webhooks) == 0 {
				cmd.Println("No Hangouts Chat Webhooks found.")
				return nil
			}

			printWebhooksTable(cmd, webhooks)

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := addClientConfigOptions(cmd, c.config); err != nil {
				return err
			}
			return c.setHangoutsChatWebhooksService(cmd, args)
		},
	}

	c.root.AddCommand(cmd)
	return addClientFlags(cmd)
}

func (c *command) setHangoutsChatWebhooksService(cmd *cobra.Command, args []string) (err error) {
	if c.hangoutsChatWebhooksService != nil {
		return nil
	}
	client, err := c.getClient(cmd)
	if err != nil {
		return err
	}
	c.hangoutsChatWebhooksService = client.HangoutsChatWebhooks
	return nil
}

type hangoutsChatWebhooksService interface {
	List(ctx context.Context) (webhooks []newreleases.Webhook, err error)
}
