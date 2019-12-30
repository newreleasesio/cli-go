// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"strconv"

	"github.com/spf13/cobra"
)

func (c *command) initProviderCmd() (err error) {
	providerCmd := &cobra.Command{
		Use:   "provider",
		Short: "Information about project providers",
	}

	optionNameAdded := "added"

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Get project providers",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			var providers []string

			added, err := cmd.Flags().GetBool(optionNameAdded)
			if err != nil {
				return err
			}

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

	listCmd.Flags().Bool(optionNameAdded, false, "get only providers for projects that are added for tracking")

	if err := addClientFlags(listCmd, c.config); err != nil {
		return err
	}
	providerCmd.AddCommand(listCmd)

	c.root.AddCommand(providerCmd)
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
	table.SetHeader([]string{"", "Name"})
	for i, name := range providers {
		table.Append([]string{strconv.Itoa(i + 1), name})
	}
	table.Render()
}
