// Copyright (c) 2019, NewReleases CLI AUTHORS.
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

func TestDiscordCmd(t *testing.T) {
	for _, tc := range []struct {
		name                   string
		discordChannelsService cmd.DiscordChannelsService
		wantOutput             string
		wantError              error
	}{
		{
			name:                   "no channels",
			discordChannelsService: newMockDiscordChannelsService(nil, nil),
			wantOutput:             "No Discord Channels found.\n",
		},
		{
			name: "with channels",
			discordChannelsService: newMockDiscordChannelsService([]newreleases.DiscordChannel{
				{
					ID:   "4qOpc9t16rpymcw7z8jwn5y6anne0sg5a9b1",
					Name: "NewReleases",
				},
				{
					ID:   "c6anne0sg9t4qOp16rpymcw7z8jwn5y5a9b1",
					Name: "Awesome project",
				},
			}, nil),
			wantOutput: "ID                                     NAME            \n4qOpc9t16rpymcw7z8jwn5y6anne0sg5a9b1   NewReleases       \nc6anne0sg9t4qOp16rpymcw7z8jwn5y5a9b1   Awesome project   \n",
		},
		{
			name:                   "error",
			discordChannelsService: newMockDiscordChannelsService(nil, errTest),
			wantError:              errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs("discord"),
				cmd.WithOutput(&outputBuf),
				cmd.WithDiscordChannelsService(tc.discordChannelsService),
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

type mockDiscordChannelsService struct {
	channels []newreleases.DiscordChannel
	err      error
}

func newMockDiscordChannelsService(channels []newreleases.DiscordChannel, err error) (s mockDiscordChannelsService) {
	return mockDiscordChannelsService{channels: channels, err: err}
}

func (s mockDiscordChannelsService) List(ctx context.Context) (channels []newreleases.DiscordChannel, err error) {
	return s.channels, s.err
}
