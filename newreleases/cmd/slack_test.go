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

func TestSlackCmd(t *testing.T) {
	for _, tc := range []struct {
		name                 string
		slackChannelsService cmd.SlackChannelsService
		wantOutput           string
		wantError            error
	}{
		{
			name:                 "no channels",
			slackChannelsService: newMockSlackChannelsService(nil, nil),
			wantOutput:           "No Slack Channels found.\n",
		},
		{
			name: "with channels",
			slackChannelsService: newMockSlackChannelsService([]newreleases.SlackChannel{
				{
					ID:       "4qOpc9t16rpymcw7z8jwn5y6anne0sg5a9b1",
					TeamName: "NewReleases",
					Channel:  "@bob",
				},
				{
					ID:       "c6anne0sg9t4qOp16rpymcw7z8jwn5y5a9b1",
					TeamName: "Awesome project",
					Channel:  "general",
				},
			}, nil),
			wantOutput: "ID                                     WORKSPACE         CHANNEL \n4qOpc9t16rpymcw7z8jwn5y6anne0sg5a9b1   NewReleases       @bob      \nc6anne0sg9t4qOp16rpymcw7z8jwn5y5a9b1   Awesome project   general   \n",
		},
		{
			name:                 "error",
			slackChannelsService: newMockSlackChannelsService(nil, errTest),
			wantError:            errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs("slack"),
				cmd.WithOutput(&outputBuf),
				cmd.WithSlackChannelsService(tc.slackChannelsService),
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

type mockSlackChannelsService struct {
	channels []newreleases.SlackChannel
	err      error
}

func newMockSlackChannelsService(channels []newreleases.SlackChannel, err error) (s mockSlackChannelsService) {
	return mockSlackChannelsService{channels: channels, err: err}
}

func (s mockSlackChannelsService) List(ctx context.Context) (channels []newreleases.SlackChannel, err error) {
	return s.channels, s.err
}
