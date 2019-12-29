// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"newreleases.io/cmd/newreleases/cmd"

	"github.com/spf13/cobra"
)

func TestMain(m *testing.M) {
	dir, err := ioutil.TempDir("", "newreleases-cmd-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	cmd.SetTestHomeDir(dir)

	os.Exit(m.Run())
}

// ExecuteT is a test function that executes command with options.
func ExecuteT(t *testing.T, opts ...Option) {
	t.Helper()

	o := &options{
		cmd:           cmd.RootCmd,
		errorRecorder: new(bytes.Buffer),
	}
	cmd.RootCmd.SetErr(o.errorRecorder)
	for _, opt := range opts {
		callback := opt.apply(o)
		if callback != nil {
			defer callback()
		}
	}
	defer cmd.NewResetCfgFileFunc()()

	if err := cmd.Execute(); err != o.wantError {
		t.Fatalf("got error %v, want %v", err, o.wantError)
	}
	if o.errorRecorder != nil {
		if errorOutput := o.errorRecorder.String(); errorOutput != "" {
			t.Fatalf("got unexpected error output:\n%q", errorOutput)
		}
	}
}

type Option interface {
	apply(*options) (callback func())
}

type options struct {
	cmd           *cobra.Command
	errorRecorder *bytes.Buffer
	wantError     error
}

func WithArgs(a ...string) Option {
	return optionFunc(func(o *options) func() {
		o.cmd.SetArgs(a)
		return nil
	})
}

func WithInput(r io.Reader) Option {
	return optionFunc(func(o *options) func() {
		o.cmd.SetIn(r)
		return nil
	})
}

func WithOutput(w io.Writer) Option {
	return optionFunc(func(o *options) func() {
		o.cmd.SetOut(w)
		return nil
	})
}

func WithErrorOutput(w io.Writer) Option {
	return optionFunc(func(o *options) func() {
		o.cmd.SetErr(w)
		o.errorRecorder = nil
		return nil
	})
}

func WithError(err error) Option {
	return optionFunc(func(o *options) func() {
		o.wantError = err
		return nil
	})
}

func WithPasswordReader(r cmd.PasswordReader) Option {
	return optionFunc(func(o *options) func() {
		orig := cmd.SetCMDPasswordReader(r)
		return func() {
			cmd.SetCMDPasswordReader(orig)
		}
	})
}

func WithAuthKeysGetter(g cmd.AuthKeysGetter) Option {
	return optionFunc(func(o *options) func() {
		orig := cmd.SetCMDAuthKeysGetter(g)
		return func() {
			cmd.SetCMDAuthKeysGetter(orig)
		}
	})
}

func WithAuthService(s cmd.AuthService) Option {
	return optionFunc(func(o *options) func() {
		orig := cmd.SetCMDAuthService(s)
		return func() {
			cmd.SetCMDAuthService(orig)
		}
	})
}

type optionFunc func(o *options) func()

func (f optionFunc) apply(o *options) func() { return f(o) }
