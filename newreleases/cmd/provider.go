// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

func (c *command) initProviderCmd() (err error) {
	optionNameAdded := "added"

	cmd := &cobra.Command{
		Use:   "provider",
		Short: "Get project providers",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			added, err := cmd.Flags().GetBool(optionNameAdded)
			if err != nil {
				return err
			}

			var providers []string
			if added {
				providers, err = c.providersService.ListAdded(ctx)
			} else {
				providers, err = c.providersService.List(ctx)
			}
			if err != nil {
				return err
			}

			if len(providers) == 0 {
				cmd.Println("No providers found.")
				return nil
			}

			printProvidersTable(cmd, providers)

			return nil
		},
		PreRunE: c.setProvidersService,
	}

	cmd.Flags().Bool(optionNameAdded, false, "get only providers for projects that are added for tracking")

	if err := addClientFlags(cmd, c.config); err != nil {
		return err
	}

	c.root.AddCommand(cmd)
	return nil
}

func (c *command) setProvidersService(cmd *cobra.Command, args []string) (err error) {
	if c.providersService != nil {
		return nil
	}
	client, err := c.getClient(cmd)
	if err != nil {
		return err
	}
	c.providersService = client.Providers
	return nil
}

type providersService interface {
	List(ctx context.Context) (providers []string, err error)
	ListAdded(ctx context.Context) (providers []string, err error)
}

func printProvidersTable(cmd *cobra.Command, providers []string) {
	table := newTable(cmd.OutOrStdout())
	table.SetHeader([]string{"ID"})
	for _, id := range providers {
		table.Append([]string{id})
	}
	table.Render()
}
