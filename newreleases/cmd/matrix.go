// Copyright (c) 2022, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"

	"github.com/spf13/cobra"
	"newreleases.io/newreleases"
)

func (c *command) initMatrixCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "matrix",
		Short: "List Matrix integrations",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			rooms, err := c.matrixRoomsService.List(ctx)
			if err != nil {
				return err
			}

			if len(rooms) == 0 {
				cmd.Println("No Matrix Rooms found.")
				return nil
			}

			printMatrixRoomsTable(cmd, rooms)

			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := addClientConfigOptions(cmd, c.config); err != nil {
				return err
			}
			return c.setMatrixRoomsService(cmd, args)
		},
	}

	c.root.AddCommand(cmd)
	return addClientFlags(cmd)
}

func (c *command) setMatrixRoomsService(cmd *cobra.Command, args []string) (err error) {
	if c.matrixRoomsService != nil {
		return nil
	}
	client, err := c.getClient(cmd)
	if err != nil {
		return err
	}
	c.matrixRoomsService = client.MatrixRooms
	return nil
}

type matrixRoomsService interface {
	List(ctx context.Context) (rooms []newreleases.MatrixRoom, err error)
}

func printMatrixRoomsTable(cmd *cobra.Command, rooms []newreleases.MatrixRoom) {
	table := newTable(cmd.OutOrStdout())
	table.SetHeader([]string{"ID", "Name", "Homeserver URL", "Internal Room ID"})
	for _, e := range rooms {
		table.Append([]string{e.ID, e.Name, e.HomeserverURL, e.InternalRoomID})
	}
	table.Render()
}
