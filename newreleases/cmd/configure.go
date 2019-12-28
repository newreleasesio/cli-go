// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bufio"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "configure",
		Short: "Provide configuration values to be stored in a file",
		Run: func(cmd *cobra.Command, args []string) {
			reader := bufio.NewReader(cmd.InOrStdin())

			authKey, err := terminalPrompt(cmd, reader, "Auth Key")
			if err != nil {
				handleError(cmd, err)
				return
			}
			if authKey == "" {
				cmd.PrintErr("No key provided.\n")
				cmd.Println("Configuration is not saved.")
				exit(1)
				return
			}

			if err := writeConfig(cmd, authKey); err != nil {
				handleError(cmd, err)
				return
			}

			cmd.Printf("Configuration saved to: %s.\n", cfgFile)
		},
	})
}
