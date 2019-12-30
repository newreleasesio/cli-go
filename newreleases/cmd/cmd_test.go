// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"newreleases.io/cmd/newreleases/cmd"
)

var homeDir string

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

	opts = append([]cmd.Option{cmd.WithHomeDir(homeDir)}, opts...)
	c, err := cmd.NewCommand(opts...)
	if err != nil {
		t.Fatal(err)
	}
	return c
}
