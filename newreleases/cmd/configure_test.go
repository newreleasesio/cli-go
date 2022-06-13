// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"newreleases.io/cmd/newreleases/cmd"
)

func TestConfigureCmd(t *testing.T) {
	for _, tc := range []struct {
		name            string
		withConfigFlag  bool
		newConfig       bool
		authKey         string
		wantOutputFunc  func(filename string) string
		wantErrorOutput string
		wantData        string
	}{
		{
			name:            "no key",
			wantOutputFunc:  func(string) string { return "Auth Key: Configuration is not saved.\n" },
			wantErrorOutput: "No key provided.\n",
		},
		{
			name:            "no key with config flag",
			withConfigFlag:  true,
			wantOutputFunc:  func(string) string { return "Auth Key: Configuration is not saved.\n" },
			wantErrorOutput: "No key provided.\n",
		},
		{
			name:    "valid key",
			authKey: "z8jwn5ne0sg5a9b4qOpc9ty6an16rpymcw71",
			wantOutputFunc: func(filename string) string {
				return fmt.Sprintf("Auth Key: Configuration saved to: %s.\n", filename)
			},
			wantData: "auth-key: z8jwn5ne0sg5a9b4qOpc9ty6an16rpymcw71\n",
		},
		{
			name:           "valid key with new config flag",
			withConfigFlag: true,
			newConfig:      true,
			authKey:        "9ty6an1z8jwn5ne0sg5a9b4qOpc6rpymcw71",
			wantOutputFunc: func(filename string) string {
				return fmt.Sprintf("Auth Key: Configuration saved to: %s.\n", filename)
			},
			wantData: "auth-key: 9ty6an1z8jwn5ne0sg5a9b4qOpc6rpymcw71\n",
		},
		{
			name:           "valid key with existing config flag",
			withConfigFlag: true,
			authKey:        "9ty6an1z8jwn5ne0sg5a9b4qOpc6rpymcw71",
			wantOutputFunc: func(filename string) string {
				return fmt.Sprintf("Auth Key: Configuration saved to: %s.\n", filename)
			},
			wantData: "auth-key: 9ty6an1z8jwn5ne0sg5a9b4qOpc6rpymcw71\n",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()

			cfgFile := filepath.Join(dir, ".newreleases.yaml")
			if !tc.newConfig {
				f, err := os.Create(cfgFile)
				if err != nil {
					t.Fatal(err)
				}
				if err := f.Close(); err != nil {
					t.Fatal(err)
				}
			}

			args := []string{"configure"}
			var setCfgFile string
			if tc.withConfigFlag {
				args = append(args, "--config", cfgFile)
			} else {
				setCfgFile = cfgFile
			}

			var outputBuf, errorOutputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs(args...),
				cmd.WithOutput(&outputBuf),
				cmd.WithErrorOutput(&errorOutputBuf),
				cmd.WithInput(strings.NewReader(tc.authKey+"\n")),
				cmd.WithCfgFile(setCfgFile),
				cmd.WithHomeDir(dir),
			).Execute(); err != nil {
				t.Fatal(err)
			}

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
				gotData, err := os.ReadFile(cfgFile)
				if err != nil {
					t.Fatal(err)
				}
				if string(gotData) != tc.wantData {
					t.Errorf("got config file data %q, want %q", string(gotData), tc.wantData)
				}
			} else {
				gotData, _ := os.ReadFile(cfgFile)
				if string(gotData) != "" {
					t.Errorf("got config file data %q, but it should not be", string(gotData))
				}
			}
		})
	}
}

func TestConfigureCmd_overwrite(t *testing.T) {
	dir := t.TempDir()

	cfgFile := filepath.Join(dir, ".newreleases.yaml")
	f, err := os.Create(cfgFile)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	testConfigre := func(t *testing.T, authKey string) {
		t.Helper()

		var outputBuf bytes.Buffer
		if err := newCommand(t,
			cmd.WithCfgFile(cfgFile),
			cmd.WithHomeDir(dir),
			cmd.WithArgs("configure"),
			cmd.WithOutput(&outputBuf),
			cmd.WithInput(strings.NewReader(authKey+"\n")),
		).Execute(); err != nil {
			t.Fatal(err)
		}

		gotOutput := outputBuf.String()
		wantOutput := fmt.Sprintf("Auth Key: Configuration saved to: %s.\n", cfgFile)
		if gotOutput != wantOutput {
			t.Errorf("got output %q, want %q", gotOutput, wantOutput)
		}

		gotData, err := os.ReadFile(cfgFile)
		if err != nil {
			t.Fatal(err)
		}
		wantData := fmt.Sprintf("auth-key: %s\n", authKey)
		if string(gotData) != wantData {
			t.Errorf("got config file data %q, want %q", string(gotData), wantData)
		}
	}

	// save first key
	testConfigre(t, "z8jwn5ne0sg5a9b4qOpc9ty6an16rpymcw71")
	// overwrite with the new key
	testConfigre(t, "9ty6an1z8jwn5ne0sg5a9b4qOpc6rpymcw71")
}
