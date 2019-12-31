// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"bytes"
	"context"
	"net"
	"testing"

	"newreleases.io/cmd/newreleases/cmd"
	"newreleases.io/newreleases"
)

func TestAuthCmd(t *testing.T) {
	_, ipNet1, err := net.ParseCIDR("127.0.0.0/8")
	if err != nil {
		t.Fatal(err)
	}
	_, ipNet2, err := net.ParseCIDR("123.33.44.1/32")
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name        string
		authService cmd.AuthService
		wantOutput  string
		wantError   error
	}{
		{
			name:        "no keys",
			authService: newMockAuthService(nil, nil),
			wantOutput:  "No auth keys found.\n",
		},
		{
			name:        "single key",
			authService: newMockAuthService([]newreleases.AuthKey{{Name: "Master", AuthorizedNetworks: []net.IPNet{*ipNet1}}}, nil),
			wantOutput:  "NAME     AUTHORIZED NETWORKS   SECRET \nMaster   127.0.0.0/8                    \n",
		},
		{
			name: "two keys",
			authService: newMockAuthService([]newreleases.AuthKey{
				{Name: "Master", AuthorizedNetworks: []net.IPNet{*ipNet1}},
				{Name: "Another", AuthorizedNetworks: []net.IPNet{*ipNet1, *ipNet2}},
			}, nil),
			wantOutput: "NAME      AUTHORIZED NETWORKS           SECRET \nMaster    127.0.0.0/8                            \nAnother   127.0.0.0/8, 123.33.44.1/32            \n",
		},
		{
			name:        "error",
			authService: newMockAuthService(nil, errTest),
			wantError:   errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs("auth"),
				cmd.WithOutput(&outputBuf),
				cmd.WithAuthService(tc.authService),
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

type mockAuthService struct {
	keys []newreleases.AuthKey
	err  error
}

func newMockAuthService(keys []newreleases.AuthKey, err error) (s mockAuthService) {
	return mockAuthService{keys: keys, err: err}
}

func (s mockAuthService) List(ctx context.Context) (keys []newreleases.AuthKey, err error) {
	return s.keys, s.err
}
