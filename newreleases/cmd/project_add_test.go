// Copyright (c) 2020, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file. {

package cmd_test

import (
	"bytes"
	"testing"

	"newreleases.io/cmd/newreleases/cmd"
)

func TestProjectCmd_Add(t *testing.T) {
	for _, tc := range []struct {
		name            string
		args            []string
		projectsService cmd.ProjectsService
		wantOutput      string
		wantError       error
	}{
		{
			name:            "minimal",
			args:            []string{"github", "golang/go"},
			projectsService: newMockProjectsService(1, nil),
			wantOutput:      "ID:         new         \nName:       golang/go   \nProvider:   github      \n",
		},
		{
			name: "full",
			args: []string{
				"github", "golang/go",
				"--email", "weekly",
				"--slack", "mdsbe60td5gwgzetyksdfeyxt4",
				"--telegram", "sdfeyxt4mdsbe60td5gwgzetyk",
				"--discord", "4mdsbe60td5gwgzetyksdfeyxt",
				"--discord", "zext4mdsbe6tyksdfey0td5gwg",
				"--hangouts-chat", "etyksdfeyxt4mdsbe60td5gwgz",
				"--microsoft-teams", "0td5gwgzextbe6tyksdfey4mds",
				"--mattermost", "wgxtzed4yksd5dfeymsbe6t0tg",
				"--rocketchat", "5dfeymsbe6t0tgwgxtzed4yksd",
				"--webhook", "tbe6tyksdfey4md0td5gwgzexs",
				"--regex-exclude", `^0\.1`,
				"--regex-exclude", `^0\.3-inverse`,
				"--exclude-prereleases",
				"--exclude-updated",
			},
			projectsService: newMockProjectsService(1, nil),
			wantOutput:      "ID:                      new                                                      \nName:                    golang/go                                                \nProvider:                github                                                   \nEmail:                   weekly                                                   \nSlack:                   mdsbe60td5gwgzetyksdfeyxt4                               \nTelegram:                sdfeyxt4mdsbe60td5gwgzetyk                               \nDiscord:                 4mdsbe60td5gwgzetyksdfeyxt, zext4mdsbe6tyksdfey0td5gwg   \nHangouts Chat:           etyksdfeyxt4mdsbe60td5gwgz                               \nMicrosoft Teams:         0td5gwgzextbe6tyksdfey4mds                               \nMattermost:              wgxtzed4yksd5dfeymsbe6t0tg                               \nRocket.Chat:             5dfeymsbe6t0tgwgxtzed4yksd                               \nWebhooks:                tbe6tyksdfey4md0td5gwgzexs                               \nRegex Exclude:           ^0\\.1                                                    \nRegex Exclude Inverse:   ^0\\.3                                                    \nExclude Pre-Releases:    yes                                                      \nExclude Updated:         yes                                                      \n",
		},
		{
			name:            "error",
			args:            []string{"github", "golang/go"},
			projectsService: newMockProjectsService(1, errTest),
			wantError:       errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs(append([]string{"project", "add"}, tc.args...)...),
				cmd.WithOutput(&outputBuf),
				cmd.WithProjectsService(tc.projectsService),
			).Execute(); err != tc.wantError {
				t.Fatalf("got error %v, want %v", err, tc.wantError)
			}

			gotOutput := outputBuf.String()
			if gotOutput != tc.wantOutput {
				t.Errorf("got output %q, want %q", gotOutput, tc.wantOutput)
			}
		})
	}
}
