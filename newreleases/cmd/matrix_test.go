// Copyright (c) 2022, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"bytes"
	"context"
	"testing"

	"newreleases.io/cmd/newreleases/cmd"
	"newreleases.io/newreleases"
)

func TestMatrixCmd(t *testing.T) {
	for _, tc := range []struct {
		name               string
		matrixRoomsService cmd.MatrixRoomsService
		wantOutput         string
		wantError          error
	}{
		{
			name:               "no rooms",
			matrixRoomsService: newMockMatrixRoomsService(nil, nil),
			wantOutput:         "No Matrix Rooms found.\n",
		},
		{
			name: "with rooms",
			matrixRoomsService: newMockMatrixRoomsService([]newreleases.MatrixRoom{
				{
					ID:             "znne04qO5y6acw7sg5a9b1pc9t16rpym8jwn",
					Name:           "NewReleases",
					HomeserverURL:  "https://matrix-client.matrix.org",
					InternalRoomID: "!CklbIcwhYKygGsulFi:matrix.org",
				},
				{
					ID:             "9t4qOp16gmcw7z8jwrpyc6anne0sn5y5a9b1",
					Name:           "Awesome project",
					HomeserverURL:  "https://matrix-client.example.com",
					InternalRoomID: "!CkGsIcwhKyulFYlbgi:example.com",
				},
			}, nil),
			wantOutput: "ID                                     NAME              HOMESERVER URL                      INTERNAL ROOM ID                \nznne04qO5y6acw7sg5a9b1pc9t16rpym8jwn   NewReleases       https://matrix-client.matrix.org    !CklbIcwhYKygGsulFi:matrix.org    \n9t4qOp16gmcw7z8jwrpyc6anne0sn5y5a9b1   Awesome project   https://matrix-client.example.com   !CkGsIcwhKyulFYlbgi:example.com   \n",
		},
		{
			name:               "error",
			matrixRoomsService: newMockMatrixRoomsService(nil, errTest),
			wantError:          errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs("matrix"),
				cmd.WithOutput(&outputBuf),
				cmd.WithMatrixRoomsService(tc.matrixRoomsService),
			).Execute(); err != tc.wantError {
				t.Fatalf("got error %v, want %v", err, tc.wantError)
			}

			gotOutput := outputBuf.String()
			if gotOutput != tc.wantOutput {
				t.Errorf("got output %q, want %q", gotOutput, tc.wantOutput)
			}
		})
	}
}

type mockMatrixRoomsService struct {
	rooms []newreleases.MatrixRoom
	err   error
}

func newMockMatrixRoomsService(rooms []newreleases.MatrixRoom, err error) mockMatrixRoomsService {
	return mockMatrixRoomsService{rooms: rooms, err: err}
}

func (s mockMatrixRoomsService) List(ctx context.Context) ([]newreleases.MatrixRoom, error) {
	return s.rooms, s.err
}
