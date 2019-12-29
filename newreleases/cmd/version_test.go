// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"bytes"
	"testing"

	nrcmd "newreleases.io/cmd"
)

func TestVersionCmd(t *testing.T) {
	var outputBuf bytes.Buffer
	ExecuteT(t, WithArgs("version"), WithOutput(&outputBuf))

	want := nrcmd.Version + "\n"
	got := outputBuf.String()
	if got != want {
		t.Errorf("got output %q, want %q", got, want)
	}
}
