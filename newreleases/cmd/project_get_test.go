// Copyright (c) 2020, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file. {

package cmd_test

import (
	"bytes"
	"testing"

	"newreleases.io/cmd/newreleases/cmd"
	"newreleases.io/newreleases"
)

func TestProjectCmd_Get(t *testing.T) {
	for _, tc := range []struct {
		name            string
		projectsService cmd.ProjectsService
		wantOutput      string
		wantError       error
	}{
		{
			name:            "no project",
			projectsService: newMockProjectsService(1, nil),
			wantOutput:      "Project not found.\n",
		},
		{
			name:            "minimal project",
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{minimalProject}),
			wantOutput:      "ID:         mdsbe60td5gwgzetyksdfeyxt4   \nName:       golang/go                    \nProvider:   github                       \n",
		},
		{
			name:            "full project",
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{fullProject}),
			wantOutput:      "ID:                      mdsbe60td5gwgzetyksdfeyxt4   \nName:                    golang/go                    \nProvider:                github                       \nEmail:                   daily                        \nSlack:                   zetyksdfeymdsbe60td5gwgxt4   \nTelegram:                sbe60td5gwgxtzetyksdfeymd4   \nDiscord:                 tyksdfeymsbegxtzed460td5gw   \nHangouts Chat:           yksdfeymsbe6t0td5gzed4wgxt   \nMicrosoft Teams:         gwgxtzed4yksdfeymsbe6t0td5   \nMattermost:              wgxtzed4yksd5dfeymsbe6t0tg   \nRocket.Chat:             5dfeymsbe6t0tgwgxtzed4yksd   \nMatrix:             4yksd5e6twgxtzdfeymsbed0tg   \nWebhooks:                e6t0td5ykgwgxtzed4eymsbsdf   \nRegex Exclude:           ^0\\.1                        \nRegex Exclude Inverse:   ^0\\.3                        \nExclude Pre-Releases:    yes                          \nExclude Updated:         yes                          \nNote:                    Initial note                 \nTags:                    33f1db7254b9                 \n",
		},
		{
			name:            "error",
			projectsService: newMockProjectsService(1, errTest),
			wantError:       errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("by id", func(t *testing.T) {
				var outputBuf bytes.Buffer
				if err := newCommand(t,
					cmd.WithArgs("project", "get", "mdsbe60td5gwgzetyksdfeyxt4"),
					cmd.WithOutput(&outputBuf),
					cmd.WithProjectsService(tc.projectsService),
				).Execute(); err != tc.wantError {
					t.Fatalf("got error %v, want %v", err, tc.wantError)
				}

				wantOutput := trimSpace(tc.wantOutput)
				gotOutput := trimSpace(outputBuf.String())
				if gotOutput != wantOutput {
					t.Errorf("got output %q, want %q", gotOutput, wantOutput)
				}
			})
			t.Run("by name", func(t *testing.T) {
				var outputBuf bytes.Buffer
				if err := newCommand(t,
					cmd.WithArgs("project", "get", "github", "golang/go"),
					cmd.WithOutput(&outputBuf),
					cmd.WithProjectsService(tc.projectsService),
				).Execute(); err != tc.wantError {
					t.Fatalf("got error %v, want %v", err, tc.wantError)
				}

				wantOutput := trimSpace(tc.wantOutput)
				gotOutput := trimSpace(outputBuf.String())
				if gotOutput != wantOutput {
					t.Errorf("got output %q, want %q", gotOutput, wantOutput)
				}
			})
		})
	}
}
