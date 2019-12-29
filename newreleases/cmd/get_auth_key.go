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

	"newreleases.io/newreleases"

	"github.com/spf13/cobra"
)

func init() {
	getAuthKeyCmd := &cobra.Command{
		Use:   "get-auth-key",
		Short: "Get API auth key and store it in the configuration",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cmd.Println("Sign in to NewReleases with your credentials")
			cmd.Println("to get available API keys and store them in local configuration file.")

			reader := bufio.NewReader(cmd.InOrStdin())

			email, err := terminalPrompt(cmd, reader, "Email")
			if err != nil {
				return err
			}
			password, err := terminalPromptPassword(cmd, "Password")
			if err != nil {
				return err
			}

			ctx, cancel := newClientContext()
			defer cancel()

			keys, err := cmdAuthKeysGetter.GetAuthKeys(ctx, email, password, newClientOptions())
			if err != nil {
				return err
			}

			count := len(keys)
			if count == 0 {
				cmd.PrintErr("No auth keys found.\n")
				cmd.Println("Go to https://newreleases.io and create an auth key.")
				return nil
			}

			var selection int
			if count > 1 {

				cmd.Println()
				printAuthKeysTable(cmd, keys)
				cmd.Println()

				for {
					in, err := terminalPrompt(cmd, reader, "Select auth key (enter row number)")
					if err != nil && err != io.EOF {
						return err
					}
					if in == "" {
						cmd.PrintErr("No key selected.\n")
						cmd.Println("Configuration is not saved.")
						return nil
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
				return err
			}
			cmd.Printf("Using auth key: %s.\n", key.Name)

			cmd.Printf("Configuration saved to: %s.\n", cfgFile)
			return nil
		},
	}

	rootCmd.AddCommand(getAuthKeyCmd)
}

var cmdAuthKeysGetter authKeysGetter = authKeysGetterFunc(newreleases.GetAuthKeys)

type authKeysGetter interface {
	GetAuthKeys(ctx context.Context, email, password string, o *newreleases.ClientOptions) (keys []newreleases.AuthKey, err error)
}

type authKeysGetterFunc func(ctx context.Context, email, password string, o *newreleases.ClientOptions) (keys []newreleases.AuthKey, err error)

func (f authKeysGetterFunc) GetAuthKeys(ctx context.Context, email, password string, o *newreleases.ClientOptions) (keys []newreleases.AuthKey, err error) {
	return f(ctx, email, password, o)
}
