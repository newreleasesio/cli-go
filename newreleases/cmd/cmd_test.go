// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"newreleases.io/cmd/newreleases/cmd"
)

var homeDir string

var errTest = errors.New("test error")

func TestMain(m *testing.M) {
	dir, err := ioutil.TempDir("", "newreleases-cmd-")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.RemoveAll(dir)

	homeDir = dir

	os.Exit(m.Run())
}

func newCommand(t *testing.T, opts ...cmd.Option) (c *cmd.Command) {
	t.Helper()

	c, err := cmd.NewCommand(append([]cmd.Option{cmd.WithHomeDir(homeDir)}, opts...)...)
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func newTime(t *testing.T, s string) (tm time.Time) {
	t.Helper()

	tm, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t.Fatal(err)
	}
	return tm
}
