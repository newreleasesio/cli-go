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

func TestTagCmd_List(t *testing.T) {
	for _, tc := range []struct {
		name        string
		tagsService cmd.TagsService
		wantOutput  string
		wantError   error
	}{
		{
			name:        "no tags",
			tagsService: newMockTagsService(nil, nil),
			wantOutput:  "No tags found.\n",
		},
		{
			name:        "with tags",
			tagsService: newMockTagsService(fullTags, nil),
			wantOutput:  "ID                                     NAME            \n33f1db7254b9   Cool       \n1d33b7254b9f   Awesome   \n",
		},
		{
			name:        "error",
			tagsService: newMockTagsService(nil, errTest),
			wantError:   errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs("tag", "list"),
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
