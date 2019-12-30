// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"bytes"
	"context"
	"errors"
	"testing"

	"newreleases.io/cmd/newreleases/cmd"
)

func TestProviderCmd(t *testing.T) {
	errTest := errors.New("test error")

	for _, tc := range []struct {
		name             string
		providersService cmd.ProvidersService
		added            bool
		wantOutput       string
		wantError        error
	}{
		{
			name:             "no providers",
			providersService: newMockProvidersService(nil, nil, nil),
			wantOutput:       "No providers found.\n",
		},
		{
			name:             "no added providers",
			added:            true,
			providersService: newMockProvidersService([]string{"github", "pypi", "npm"}, nil, nil),
			wantOutput:       "No providers found.\n",
		},
		{
			name:             "providers",
			providersService: newMockProvidersService([]string{"github", "pypi", "cargo", "dockerhub"}, []string{"github", "pypi"}, nil),
			wantOutput:       "|   |   NAME    |\n|---|-----------|\n| 1 | github    |\n| 2 | pypi      |\n| 3 | cargo     |\n| 4 | dockerhub |\n",
		},
		{
			name:             "added providers",
			added:            true,
			providersService: newMockProvidersService([]string{"github", "pypi", "yarn", "dockerhub"}, []string{"github", "pypi"}, nil),
			wantOutput:       "|   |  NAME  |\n|---|--------|\n| 1 | github |\n| 2 | pypi   |\n",
		},
		{
			name:             "error",
			providersService: newMockProvidersService(nil, nil, errTest),
			wantError:        errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			args := []string{"provider", "list"}
			if tc.added {
				args = append(args, "--added")
			}
			var outputBuf bytes.Buffer
			c := newCommand(t,
				cmd.WithArgs(args...),
				cmd.WithOutput(&outputBuf),
				cmd.WithProvidersService(tc.providersService),
			)
			if err := c.Execute(); err != tc.wantError {
				t.Fatalf("got error %v, want %v", err, tc.wantError)
			}

			gotOutput := outputBuf.String()
			if gotOutput != tc.wantOutput {
				t.Errorf("got output %q, want %q", gotOutput, tc.wantOutput)
			}
		})
	}
}

type mockProvidersService struct {
	providers      []string
	addedProviders []string
	err            error
}

func newMockProvidersService(providers, addedProviders []string, err error) (s mockProvidersService) {
	return mockProvidersService{providers: providers, addedProviders: addedProviders, err: err}
}

func (s mockProvidersService) List(ctx context.Context) (providers []string, err error) {
	return s.providers, s.err
}

func (s mockProvidersService) ListAdded(ctx context.Context) (providers []string, err error) {
	return s.addedProviders, s.err
}
