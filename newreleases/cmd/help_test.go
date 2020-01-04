// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"bytes"
	"strings"
	"testing"

	"newreleases.io/cmd/newreleases/cmd"
)

func TestRootCmdHelp(t *testing.T) {
	for _, arg := range []string{
		"",
		"-h",
		"--help",
	} {
		var outputBuf bytes.Buffer
		if err := newCommand(t,
			cmd.WithArgs(arg),
			cmd.WithOutput(&outputBuf),
		).Execute(); err != nil {
			t.Fatal(err)
		}

		want := "NewReleases is a release tracker for software engineers"
		got := outputBuf.String()
		if !strings.Contains(got, want) {
			t.Errorf("output %q for arg %q, does not contain %q", got, arg, want)
		}
	}
}
