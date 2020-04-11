// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
	"newreleases.io/newreleases"
)

func (c *command) initAuthCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Get API authentication keys",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			keys, err := c.authService.List(ctx)
			if err != nil {
				return err
			}

			if len(keys) == 0 {
				cmd.Println("No auth keys found.")
				return nil
			}

			printAuthKeysTable(cmd, keys)

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := addClientConfigOptions(cmd, c.config); err != nil {
				return err
			}
			return c.setAuthService(cmd, args)
		},
	}

	c.root.AddCommand(cmd)
	return addClientFlags(cmd)
}

func (c *command) setAuthService(cmd *cobra.Command, args []string) (err error) {
	if c.authService != nil {
		return nil
	}
	client, err := c.getClient(cmd)
	if err != nil {
		return err
	}
	c.authService = client.Auth
	return nil
}

type authService interface {
	List(ctx context.Context) (keys []newreleases.AuthKey, err error)
}

func printAuthKeysTable(cmd *cobra.Command, keys []newreleases.AuthKey) {
	table := newTable(cmd.OutOrStdout())
	table.SetHeader([]string{"Name", "Authorized Networks", "Secret"})
	for _, key := range keys {
		var authorizedNetworks []string
		for _, an := range key.AuthorizedNetworks {
			authorizedNetworks = append(authorizedNetworks, an.String())
		}
		table.Append([]string{key.Name, strings.Join(authorizedNetworks, ", "), key.Secret})
	}
	table.Render()
}
