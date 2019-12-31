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

func TestMicrosoftTeamsCmd(t *testing.T) {
	for _, tc := range []struct {
		name                          string
		microsoftTeamsWebhooksService cmd.MicrosoftTeamsWebhooksService
		wantOutput                    string
		wantError                     error
	}{
		{
			name:                          "no webhooks",
			microsoftTeamsWebhooksService: newMockMicrosoftTeamsWebhooksService(nil, nil),
			wantOutput:                    "No Microsoft Teams Webhooks found.\n",
		},
		{
			name: "with webhooks",
			microsoftTeamsWebhooksService: newMockMicrosoftTeamsWebhooksService([]newreleases.Webhook{
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
			name:                          "error",
			microsoftTeamsWebhooksService: newMockMicrosoftTeamsWebhooksService(nil, errTest),
			wantError:                     errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs("microsoft-teams"),
				cmd.WithOutput(&outputBuf),
				cmd.WithMicrosoftTeamsWebhooksService(tc.microsoftTeamsWebhooksService),
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

type mockMicrosoftTeamsWebhooksService struct {
	webhooks []newreleases.Webhook
	err      error
}

func newMockMicrosoftTeamsWebhooksService(webhooks []newreleases.Webhook, err error) (s mockMicrosoftTeamsWebhooksService) {
	return mockMicrosoftTeamsWebhooksService{webhooks: webhooks, err: err}
}

func (s mockMicrosoftTeamsWebhooksService) List(ctx context.Context) (webhooks []newreleases.Webhook, err error) {
	return s.webhooks, s.err
}
