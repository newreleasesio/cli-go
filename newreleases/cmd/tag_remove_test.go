// Copyright (c) 2022, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"bytes"
	"testing"

	"newreleases.io/cmd/newreleases/cmd"
)

func TestTagCmd_Remove(t *testing.T) {
	for _, tc := range []struct {
		name        string
		tagsService cmd.TagsService
		wantOutput  string
		wantError   error
	}{
		{
			name:        "no tags",
			tagsService: newMockTagsService(nil, nil),
			wantOutput:  "Tag not found.\n",
		},
		{
			name:        "with tags",
			tagsService: newMockTagsService(fullTags, nil),
			wantOutput:  "",
		},
		{
			name:        "error",
			tagsService: newMockTagsService(fullTags, errTest),
			wantError:   errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs("tag", "remove", "33f1db7254b9"),
				cmd.WithOutput(&outputBuf),
				cmd.WithTagsService(tc.tagsService),
			).Execute(); err != tc.wantError {
				t.Fatalf("got error %v, want %v", err, tc.wantError)
			}

			wantOutput := trimSpace(tc.wantOutput)
			gotOutput := trimSpace(outputBuf.String())
			if gotOutput != wantOutput {
				t.Errorf("got output %q, want %q", gotOutput, wantOutput)
			}
		})
	}
}
