// Copyright (c) 2021, NewReleases CLI AUTHORS.
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

func TestRocketchatCmd(t *testing.T) {
	for _, tc := range []struct {
		name                      string
		rocketchatWebhooksService cmd.RocketchatWebhooksService
		wantOutput                string
		wantError                 error
	}{
		{
			name:                      "no webhooks",
			rocketchatWebhooksService: newMockRocketchatWebhooksService(nil, nil),
			wantOutput:                "No Rocket.Chat Webhooks found.\n",
		},
		{
			name: "with webhooks",
			rocketchatWebhooksService: newMockRocketchatWebhooksService([]newreleases.Webhook{
				{
					ID:   "abOpc9t16rpymcw7z8jwn5y6anne0sg5a9b1",
					Name: "NewReleases",
				},
				{
					ID:   "f1anne0sg9t4qOp16rpymcw7z8jwn5y5a9b1",
					Name: "Awesome project",
				},
			}, nil),
			wantOutput: "ID                                     NAME            \nabOpc9t16rpymcw7z8jwn5y6anne0sg5a9b1   NewReleases       \nf1anne0sg9t4qOp16rpymcw7z8jwn5y5a9b1   Awesome project   \n",
		},
		{
			name:                      "error",
			rocketchatWebhooksService: newMockRocketchatWebhooksService(nil, errTest),
			wantError:                 errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs("rocketchat"),
				cmd.WithOutput(&outputBuf),
				cmd.WithRocketchatWebhooksService(tc.rocketchatWebhooksService),
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

type mockRocketchatWebhooksService struct {
	webhooks []newreleases.Webhook
	err      error
}

func newMockRocketchatWebhooksService(webhooks []newreleases.Webhook, err error) (s mockRocketchatWebhooksService) {
	return mockRocketchatWebhooksService{webhooks: webhooks, err: err}
}

func (s mockRocketchatWebhooksService) List(ctx context.Context) (webhooks []newreleases.Webhook, err error) {
	return s.webhooks, s.err
}
