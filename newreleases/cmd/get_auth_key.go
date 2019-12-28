// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bufio"
	"context"
	"io"
	"strconv"
	"time"

	"newreleases.io/newreleases"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "get-auth-key",
		Short: "Get API auth key and store it in the configuration",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println("Sign in to NewReleases with your credentials")
			cmd.Println("to get available API keys and store them in local configuration file.")

			reader := bufio.NewReader(cmd.InOrStdin())

			email, err := terminalPrompt(cmd, reader, "Email")
			if err != nil {
				handleError(cmd, err)
				return
			}
			password, err := terminalPromptPassword(cmd, "Password")
			if err != nil {
				handleError(cmd, err)
				return
			}

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			keys, err := cmdAuthKeysGetter.GetAuthKeys(ctx, email, password, nil)
			if err != nil {
				handleError(cmd, err)
			}

			count := len(keys)
			if count == 0 {
				cmd.PrintErr("No auth keys found.\n")
				cmd.Println("Go to https://newreleases.io and create an auth key.")
				exit(1)
				return
			}

			var selection int
			if count > 1 {

				table := tablewriter.NewWriter(cmd.OutOrStdout())
				table.SetBorder(false)
				table.SetHeader([]string{"", "Name"})
				for i, key := range keys {
					table.Append([]string{strconv.Itoa(i + 1), key.Name})
				}
				cmd.Println()
				table.Render()
				cmd.Println()

				for {
					in, err := terminalPrompt(cmd, reader, "Select auth key (enter row number)")
					if err != nil && err != io.EOF {
						handleError(cmd, err)
						return
					}
					if in == "" {
						cmd.PrintErr("No key selected.\n")
						cmd.Println("Configuration is not saved.")
						exit(1)
						return
					}

					i, err := strconv.Atoi(in)
					if err != nil || i <= 0 || i > count {
						cmd.PrintErr("Invalid row number.\n")
						continue
					}
					selection = i - 1
					break
				}
			}

			key := keys[selection]
			if err := writeConfig(cmd, key.Secret); err != nil {
				handleError(cmd, err)
				return
			}
			cmd.Printf("Using auth key: %s.\n", key.Name)

			cmd.Printf("Configuration saved to: %s.\n", cfgFile)
		},
	})
}

var cmdAuthKeysGetter authKeysGetter = authKeysGetterFunc(newreleases.GetAuthKeys)

type authKeysGetter interface {
	GetAuthKeys(ctx context.Context, email, password string, o *newreleases.ClientOptions) (keys []newreleases.AuthKey, err error)
}

type authKeysGetterFunc func(ctx context.Context, email, password string, o *newreleases.ClientOptions) (keys []newreleases.AuthKey, err error)

func (f authKeysGetterFunc) GetAuthKeys(ctx context.Context, email, password string, o *newreleases.ClientOptions) (keys []newreleases.AuthKey, err error) {
	return f(ctx, email, password, o)
}
