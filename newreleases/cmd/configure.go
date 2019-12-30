// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"bufio"

	"github.com/spf13/cobra"
)

func (c *command) initConfigureCmd() {
	c.root.AddCommand(&cobra.Command{
		Use:   "configure",
		Short: "Provide configuration values to be stored in a file",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			reader := bufio.NewReader(cmd.InOrStdin())

			authKey, err := terminalPrompt(cmd, reader, "Auth Key")
			if err != nil {
				return err
			}
			if authKey == "" {
				cmd.PrintErr("No key provided.\n")
				cmd.Println("Configuration is not saved.")
				return nil
			}

			if err := c.writeConfig(cmd, authKey); err != nil {
				return err
			}

			cmd.Printf("Configuration saved to: %s.\n", c.cfgFile)
			return nil
		},
	})
}
