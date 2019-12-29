// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootCmdHelp(t *testing.T) {
	for _, arg := range []string{
		"",
		"-h",
		"--help",
	} {
		var outputBuf bytes.Buffer
		ExecuteT(t, WithArgs(arg), WithOutput(&outputBuf))

		want := "Release tracker for software engineers"
		got := outputBuf.String()
		if !strings.Contains(got, want) {
			t.Errorf("output %q for arg %q, does not contain %q", got, arg, want)
		}
	}
}
