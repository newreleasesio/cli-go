// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"newreleases.io/newreleases"
)

func init() {
	authCmd := &cobra.Command{
		Use:   "auth",
		Short: "Information about API authentication",
	}

	listCmd := &cobra.Command{
		Use:   "list",
		Short: "Get all API authentication keys",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext()
			defer cancel()

			keys, err := cmdAuthService.List(ctx)
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
		PreRunE: setCmdAuthService,
	}

	addClientFlags(listCmd)

	authCmd.AddCommand(listCmd)

	rootCmd.AddCommand(authCmd)
}

var cmdAuthService authService

func setCmdAuthService(cmd *cobra.Command, args []string) (err error) {
	if cmdAuthService != nil {
		return nil
	}
	client, err := newClient()
	if err != nil {
		return err
	}
	cmdAuthService = client.Auth
	return nil
}

type authService interface {
	List(ctx context.Context) (keys []newreleases.AuthKey, err error)
}

func printAuthKeysTable(cmd *cobra.Command, keys []newreleases.AuthKey) {
	table := tablewriter.NewWriter(cmd.OutOrStdout())
	table.SetBorder(false)
	table.SetHeader([]string{"", "Name", "Authorized Networks"})
	for i, key := range keys {
		var authorizedNetworks []string
		for _, an := range key.AuthorizedNetworks {
			authorizedNetworks = append(authorizedNetworks, an.String())
		}
		table.Append([]string{strconv.Itoa(i + 1), key.Name, strings.Join(authorizedNetworks, ", ")})
	}
	table.Render()
}
