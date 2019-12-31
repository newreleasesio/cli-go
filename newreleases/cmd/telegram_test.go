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

func TestTelegramCmd(t *testing.T) {
	for _, tc := range []struct {
		name                 string
		telegramChatsService cmd.TelegramChatsService
		wantOutput           string
		wantError            error
	}{
		{
			name:                 "no chats",
			telegramChatsService: newMockTelegramChatssService(nil, nil),
			wantOutput:           "No Telegram Chats found.\n",
		},
		{
			name: "with chats",
			telegramChatsService: newMockTelegramChatssService([]newreleases.TelegramChat{
				{
					ID:   "4qOpc9t16rpymcw7z8jwn5y6anne0sg5a9b1",
					Name: "NewReleases",
				},
				{
					ID:   "c6anne0sg9t4qOp16rpymcw7z8jwn5y5a9b1",
					Name: "Awesome project",
				},
			}, nil),
			wantOutput: "ID                                     CHAT              TYPE \n4qOpc9t16rpymcw7z8jwn5y6anne0sg5a9b1   NewReleases              \nc6anne0sg9t4qOp16rpymcw7z8jwn5y5a9b1   Awesome project          \n",
		},
		{
			name:                 "error",
			telegramChatsService: newMockTelegramChatssService(nil, errTest),
			wantError:            errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs("telegram"),
				cmd.WithOutput(&outputBuf),
				cmd.WithTelegramChatsService(tc.telegramChatsService),
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

type mockTelegramChatsService struct {
	chats []newreleases.TelegramChat
	err   error
}

func newMockTelegramChatssService(chats []newreleases.TelegramChat, err error) (s mockTelegramChatsService) {
	return mockTelegramChatsService{chats: chats, err: err}
}

func (s mockTelegramChatsService) List(ctx context.Context) (chats []newreleases.TelegramChat, err error) {
	return s.chats, s.err
}
