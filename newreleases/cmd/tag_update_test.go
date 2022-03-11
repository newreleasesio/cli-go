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

func TestTagCmd_Update(t *testing.T) {
	for _, tc := range []struct {
		name        string
		tagsService cmd.TagsService
		args        []string
		wantOutput  string
		wantError   error
	}{
		{
			name:        "name",
			tagsService: newMockTagsService(fullTags, nil),
			args:        []string{"--name", "Interesting"},
			wantOutput:  "ID: 33f1db7254b9\nName:   Interesting\n",
		},
		{
			name:        "no name flag",
			tagsService: newMockTagsService(fullTags, nil),
			wantOutput:  "Option --name is required.\n",
		},
		{
			name:        "empty name",
			tagsService: newMockTagsService(fullTags, nil),
			args:        []string{"--name", ""},
			wantOutput:  "Option --name is required.\n",
		},
		{
			name:        "error",
			tagsService: newMockTagsService(fullTags, errTest),
			args:        []string{"--name", "Interesting"},
			wantError:   errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs(append([]string{"tag", "update", "33f1db7254b9"}, tc.args...)...),
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
