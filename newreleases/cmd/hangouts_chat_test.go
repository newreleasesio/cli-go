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

func TestHangoutsChatCmd(t *testing.T) {
	for _, tc := range []struct {
		name                        string
		hangoutsChatWebhooksService cmd.HangoutsChatWebhooksService
		wantOutput                  string
		wantError                   error
	}{
		{
			name:                        "no webhooks",
			hangoutsChatWebhooksService: newMockHangoutsChatWebhooksService(nil, nil),
			wantOutput:                  "No Hangouts Chat Webhooks found.\n",
		},
		{
			name: "with webhooks",
			hangoutsChatWebhooksService: newMockHangoutsChatWebhooksService([]newreleases.Webhook{
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
			name:                        "error",
			hangoutsChatWebhooksService: newMockHangoutsChatWebhooksService(nil, errTest),
			wantError:                   errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs("hangouts-chat"),
				cmd.WithOutput(&outputBuf),
				cmd.WithHangoutsChatWebhooksService(tc.hangoutsChatWebhooksService),
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

type mockHangoutsChatWebhooksService struct {
	webhooks []newreleases.Webhook
	err      error
}

func newMockHangoutsChatWebhooksService(webhooks []newreleases.Webhook, err error) (s mockHangoutsChatWebhooksService) {
	return mockHangoutsChatWebhooksService{webhooks: webhooks, err: err}
}

func (s mockHangoutsChatWebhooksService) List(ctx context.Context) (webhooks []newreleases.Webhook, err error) {
	return s.webhooks, s.err
}
