// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"newreleases.io/cmd/newreleases/cmd"
	"newreleases.io/newreleases"
)

func TestGetAuthKeyCmd(t *testing.T) {
	for _, tc := range []struct {
		name            string
		withConfigFlag  bool
		input           string
		authKeysGetter  cmd.AuthKeysGetter
		wantOutputFunc  func(filename string) string
		wantErrorOutput string
		wantData        string
		wantError       error
	}{
		{
			name:           "empty input",
			input:          "\n",
			authKeysGetter: newMockAuthKeysGetter("", "myPassword", nil, nil),
			wantOutputFunc: func(string) string {
				return "Sign in to NewReleases with your credentials\nto get available API keys and store them in local configuration file.\nEmail: Password: \nGo to https://newreleases.io and create an auth key.\n"
			},
			wantErrorOutput: "No auth keys found.\n",
		},
		{
			name:           "unauthorized",
			input:          "me@newreleases.io\n",
			authKeysGetter: newMockAuthKeysGetter("me@newreleases.io", "wrongPassword", nil, nil),
			wantOutputFunc: func(string) string {
				return "Sign in to NewReleases with your credentials\nto get available API keys and store them in local configuration file.\nEmail: Password: \n"
			},
			wantError: newreleases.ErrUnauthorized,
		},
		{
			name:           "single key",
			input:          "me@newreleases.io\n",
			authKeysGetter: newMockAuthKeysGetter("me@newreleases.io", "myPassword", []newreleases.AuthKey{{Name: "Master", Secret: "z8jwn5ne0sg5a9b4qOpc9ty6an16rpymcw71"}}, nil),
			wantOutputFunc: func(filename string) string {
				return fmt.Sprintf("Sign in to NewReleases with your credentials\nto get available API keys and store them in local configuration file.\nEmail: Password: \nUsing auth key: Master.\nConfiguration saved to: %s.\n", filename)
			},
			wantData: "api-endpoint: \"\"\nauth-key: z8jwn5ne0sg5a9b4qOpc9ty6an16rpymcw71\ntimeout: 30s\n",
		},
		{
			name:           "single key with config flag",
			withConfigFlag: true,
			input:          "me@newreleases.io\n",
			authKeysGetter: newMockAuthKeysGetter("me@newreleases.io", "myPassword", []newreleases.AuthKey{{Name: "Master", Secret: "z8jwn5ne0sg5a9b4qOpc9ty6an16rpymcw71"}}, nil),
			wantOutputFunc: func(filename string) string {
				return fmt.Sprintf("Sign in to NewReleases with your credentials\nto get available API keys and store them in local configuration file.\nEmail: Password: \nUsing auth key: Master.\nConfiguration saved to: %s.\n", filename)
			},
			wantData: "api-endpoint: \"\"\nauth-key: z8jwn5ne0sg5a9b4qOpc9ty6an16rpymcw71\ntimeout: 30s\n",
		},
		{
			name:  "multiple keys select first",
			input: "me@newreleases.io\n1\n",
			authKeysGetter: newMockAuthKeysGetter("me@newreleases.io", "myPassword", []newreleases.AuthKey{
				{Name: "Master", Secret: "z8jwn5ne0sg5a9b4qOpc9ty6an16rpymcw71"},
				{Name: "Secondary", Secret: "ne0sg5a9b4qOpc9ty6az8jwn5n16rpymcw71"},
			}, nil),
			wantOutputFunc: func(filename string) string {
				return fmt.Sprintf("Sign in to NewReleases with your credentials\nto get available API keys and store them in local configuration file.\nEmail: Password: \n\n    |   NAME    | AUTHORIZED NETWORKS  \n----+-----------+----------------------\n  1 | Master    |                      \n  2 | Secondary |                      \n\nSelect auth key (enter row number): Using auth key: Master.\nConfiguration saved to: %s.\n", filename)
			},
			wantData: "api-endpoint: \"\"\nauth-key: z8jwn5ne0sg5a9b4qOpc9ty6an16rpymcw71\ntimeout: 30s\n",
		},
		{
			name:  "multiple keys select second",
			input: "me@newreleases.io\n2\n",
			authKeysGetter: newMockAuthKeysGetter("me@newreleases.io", "myPassword", []newreleases.AuthKey{
				{Name: "Master", Secret: "z8jwn5ne0sg5a9b4qOpc9ty6an16rpymcw71"},
				{Name: "Secondary", Secret: "ne0sg5a9b4qOpc9ty6az8jwn5n16rpymcw71"},
			}, nil),
			wantOutputFunc: func(filename string) string {
				return fmt.Sprintf("Sign in to NewReleases with your credentials\nto get available API keys and store them in local configuration file.\nEmail: Password: \n\n    |   NAME    | AUTHORIZED NETWORKS  \n----+-----------+----------------------\n  1 | Master    |                      \n  2 | Secondary |                      \n\nSelect auth key (enter row number): Using auth key: Secondary.\nConfiguration saved to: %s.\n", filename)
			},
			wantData: "api-endpoint: \"\"\nauth-key: ne0sg5a9b4qOpc9ty6az8jwn5n16rpymcw71\ntimeout: 30s\n",
		},
		{
			name:           "multiple keys select second with config flag",
			withConfigFlag: true,
			input:          "me@newreleases.io\n2\n",
			authKeysGetter: newMockAuthKeysGetter("me@newreleases.io", "myPassword", []newreleases.AuthKey{
				{Name: "Master", Secret: "z8jwn5ne0sg5a9b4qOpc9ty6an16rpymcw71"},
				{Name: "Secondary", Secret: "ne0sg5a9b4qOpc9ty6az8jwn5n16rpymcw71"},
			}, nil),
			wantOutputFunc: func(filename string) string {
				return fmt.Sprintf("Sign in to NewReleases with your credentials\nto get available API keys and store them in local configuration file.\nEmail: Password: \n\n    |   NAME    | AUTHORIZED NETWORKS  \n----+-----------+----------------------\n  1 | Master    |                      \n  2 | Secondary |                      \n\nSelect auth key (enter row number): Using auth key: Secondary.\nConfiguration saved to: %s.\n", filename)
			},
			wantData: "api-endpoint: \"\"\nauth-key: ne0sg5a9b4qOpc9ty6az8jwn5n16rpymcw71\ntimeout: 30s\n",
		},
		{
			name:  "multiple keys select none",
			input: "me@newreleases.io\n\n",
			authKeysGetter: newMockAuthKeysGetter("me@newreleases.io", "myPassword", []newreleases.AuthKey{
				{Name: "Master", Secret: "z8jwn5ne0sg5a9b4qOpc9ty6an16rpymcw71"},
				{Name: "Secondary", Secret: "ne0sg5a9b4qOpc9ty6az8jwn5n16rpymcw71"},
			}, nil),
			wantOutputFunc: func(string) string {
				return "Sign in to NewReleases with your credentials\nto get available API keys and store them in local configuration file.\nEmail: Password: \n\n    |   NAME    | AUTHORIZED NETWORKS  \n----+-----------+----------------------\n  1 | Master    |                      \n  2 | Secondary |                      \n\nSelect auth key (enter row number): Configuration is not saved.\n"
			},
			wantErrorOutput: "No key selected.\n",
		},
		{
			name:  "multiple keys select invalid then first",
			input: "me@newreleases.io\naaa\n1\n",
			authKeysGetter: newMockAuthKeysGetter("me@newreleases.io", "myPassword", []newreleases.AuthKey{
				{Name: "Master", Secret: "z8jwn5ne0sg5a9b4qOpc9ty6an16rpymcw71"},
				{Name: "Secondary", Secret: "ne0sg5a9b4qOpc9ty6az8jwn5n16rpymcw71"},
			}, nil),
			wantOutputFunc: func(filename string) string {
				return fmt.Sprintf("Sign in to NewReleases with your credentials\nto get available API keys and store them in local configuration file.\nEmail: Password: \n\n    |   NAME    | AUTHORIZED NETWORKS  \n----+-----------+----------------------\n  1 | Master    |                      \n  2 | Secondary |                      \n\nSelect auth key (enter row number): Select auth key (enter row number): Using auth key: Master.\nConfiguration saved to: %s.\n", filename)
			},
			wantErrorOutput: "Invalid row number.\n",
			wantData:        "api-endpoint: \"\"\nauth-key: z8jwn5ne0sg5a9b4qOpc9ty6an16rpymcw71\ntimeout: 30s\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir, err := ioutil.TempDir("", "newreleases-cmd-")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(dir)

			cfgFile := filepath.Join(dir, ".newreleases.yaml")
			f, err := os.Create(cfgFile)
			if err != nil {
				t.Fatal(err)
			}
			if err := f.Close(); err != nil {
				t.Fatal(err)
			}

			args := []string{"get-auth-key"}
			if tc.withConfigFlag {
				args = append(args, "--config", cfgFile)
			} else {
				defer cmd.SetCfgFile(cfgFile)()
			}

			var outputBuf, errorOutputBuf bytes.Buffer
			cmd.ExecuteT(t,
				cmd.WithArgs(args...),
				cmd.WithOutput(&outputBuf),
				cmd.WithErrorOutput(&errorOutputBuf),
				cmd.WithInput(strings.NewReader(tc.input)),
				cmd.WithError(tc.wantError),
				cmd.WithPasswordReader(newMockPasswordReader("myPassword", nil)),
				cmd.WithAuthKeysGetter(tc.authKeysGetter),
			)

			gotOutput := outputBuf.String()
			if wantOutput := tc.wantOutputFunc(cfgFile); wantOutput != "" {
				wantOutput := wantOutput
				if gotOutput != wantOutput {
					t.Errorf("got output %q, want %q", gotOutput, wantOutput)
				}
			} else {
				if gotOutput != "" {
					t.Errorf("got output %q, but it should not be", gotOutput)
				}
			}

			gotErrorOutput := errorOutputBuf.String()
			if gotErrorOutput != tc.wantErrorOutput {
				t.Errorf("got error output %q, want %q", gotErrorOutput, tc.wantErrorOutput)
			}

			if tc.wantData != "" {
				gotData, err := ioutil.ReadFile(cfgFile)
				if err != nil {
					t.Fatal(err)
				}
				if string(gotData) != tc.wantData {
					t.Errorf("got config file data %q, want %q", string(gotData), tc.wantData)
				}
			} else {
				gotData, _ := ioutil.ReadFile(cfgFile)
				if string(gotData) != "" {
					t.Errorf("got config file data %q, but it should not be", string(gotData))
				}
			}
		})
	}
}

type mockPasswordReader struct {
	password string
	err      error
}

func newMockPasswordReader(password string, err error) (r mockPasswordReader) {
	return mockPasswordReader{
		password: password,
		err:      err,
	}
}

func (r mockPasswordReader) ReadPassword() (password string, err error) {
	return r.password, r.err
}

type mockAuthKeysGetter struct {
	email    string
	password string
	keys     []newreleases.AuthKey
	err      error
}

func newMockAuthKeysGetter(email, password string, keys []newreleases.AuthKey, err error) (g mockAuthKeysGetter) {
	return mockAuthKeysGetter{
		email:    email,
		password: password,
		keys:     keys,
		err:      err,
	}
}

func (g mockAuthKeysGetter) GetAuthKeys(_ context.Context, email, password string, _ *newreleases.ClientOptions) (keys []newreleases.AuthKey, err error) {
	if email != g.email || password != g.password {
		return nil, newreleases.ErrUnauthorized
	}
	return g.keys, g.err
}
