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

func (c *command) initTelegramCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "telegram",
		Short: "List Telegram integrations",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			chats, err := c.telegramChatsService.List(ctx)
			if err != nil {
				return err
			}

			if len(chats) == 0 {
				cmd.Println("No Telegram Chats found.")
				return nil
			}

			printTelegramChatsTable(cmd, chats)

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := addClientFlags(cmd, c.config); err != nil {
				return err
			}
			return c.setTelegramChatsService(cmd, args)
		},
	}

	c.root.AddCommand(cmd)
	return nil
}

func (c *command) setTelegramChatsService(cmd *cobra.Command, args []string) (err error) {
	if c.telegramChatsService != nil {
		return nil
	}
	client, err := c.getClient(cmd)
	if err != nil {
		return err
	}
	c.telegramChatsService = client.TelegramChats
	return nil
}

type telegramChatsService interface {
	List(ctx context.Context) (chats []newreleases.TelegramChat, err error)
}

func printTelegramChatsTable(cmd *cobra.Command, chats []newreleases.TelegramChat) {
	table := newTable(cmd.OutOrStdout())
	table.SetHeader([]string{"ID", "Chat", "Type"})
	for _, e := range chats {
		table.Append([]string{e.ID, e.Name, e.Type})
	}
	table.Render()
}
