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

func TestProjectCmd_Update(t *testing.T) {
	for _, tc := range []struct {
		name            string
		args            []string
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
			name:            "no change",
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{minimalProject}),
			wantOutput:      "ID:         mdsbe60td5gwgzetyksdfeyxt4   \nName:       golang/go                    \nProvider:   github                       \n",
		},
		{
			name: "set email",
			args: []string{
				"--email", "weekly",
			},
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{minimalProject}),
			wantOutput:      "ID:         mdsbe60td5gwgzetyksdfeyxt4   \nName:       golang/go                    \nProvider:   github                       \nEmail:      weekly                       \n",
		},
		{
			name: "update email",
			args: []string{
				"--email", "weekly",
			},
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{fullProject}),
			wantOutput:      "ID:                      mdsbe60td5gwgzetyksdfeyxt4   \nName:                    golang/go                    \nProvider:                github                       \nEmail:                   weekly                       \nSlack:                   zetyksdfeymdsbe60td5gwgxt4   \nTelegram:                sbe60td5gwgxtzetyksdfeymd4   \nDiscord:                 tyksdfeymsbegxtzed460td5gw   \nHangouts Chat:           yksdfeymsbe6t0td5gzed4wgxt   \nMicrosoft Teams:         gwgxtzed4yksdfeymsbe6t0td5   \nMattermost:              wgxtzed4yksd5dfeymsbe6t0tg   \nWebhooks:                e6t0td5ykgwgxtzed4eymsbsdf   \nRegex Exclude:           ^0\\.1                        \nRegex Exclude Inverse:   ^0\\.3                        \nExclude Pre-Releases:    yes                          \nExclude Updated:         yes                          \n",
		},
		{
			name: "update slack",
			args: []string{
				"--slack", "ymdsbe60td5gwgxt4zetyksdfe",
				"--slack", "gwgxt4zetyksdfeymdsbe60td5",
			},
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{fullProject}),
			wantOutput:      "ID:                      mdsbe60td5gwgzetyksdfeyxt4                               \nName:                    golang/go                                                \nProvider:                github                                                   \nSlack:                   ymdsbe60td5gwgxt4zetyksdfe, gwgxt4zetyksdfeymdsbe60td5   \nTelegram:                sbe60td5gwgxtzetyksdfeymd4                               \nDiscord:                 tyksdfeymsbegxtzed460td5gw                               \nHangouts Chat:           yksdfeymsbe6t0td5gzed4wgxt                               \nMicrosoft Teams:         gwgxtzed4yksdfeymsbe6t0td5                               \nMattermost:              wgxtzed4yksd5dfeymsbe6t0tg                               \nWebhooks:                e6t0td5ykgwgxtzed4eymsbsdf                               \nRegex Exclude:           ^0\\.1                                                    \nRegex Exclude Inverse:   ^0\\.3                                                    \nExclude Pre-Releases:    yes                                                      \nExclude Updated:         yes                                                      \n",
		},
		{
			name: "remove slack",
			args: []string{
				"--slack-remove",
			},
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{fullProject}),
			wantOutput:      "ID:                      mdsbe60td5gwgzetyksdfeyxt4   \nName:                    golang/go                    \nProvider:                github                       \nTelegram:                sbe60td5gwgxtzetyksdfeymd4   \nDiscord:                 tyksdfeymsbegxtzed460td5gw   \nHangouts Chat:           yksdfeymsbe6t0td5gzed4wgxt   \nMicrosoft Teams:         gwgxtzed4yksdfeymsbe6t0td5   \nMattermost:              wgxtzed4yksd5dfeymsbe6t0tg   \nWebhooks:                e6t0td5ykgwgxtzed4eymsbsdf   \nRegex Exclude:           ^0\\.1                        \nRegex Exclude Inverse:   ^0\\.3                        \nExclude Pre-Releases:    yes                          \nExclude Updated:         yes                          \n",
		},
		{
			name: "include prereleases",
			args: []string{
				"--exclude-prereleases=false",
			},
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{fullProject}),
			wantOutput:      "ID:                      mdsbe60td5gwgzetyksdfeyxt4   \nName:                    golang/go                    \nProvider:                github                       \nSlack:                   zetyksdfeymdsbe60td5gwgxt4   \nTelegram:                sbe60td5gwgxtzetyksdfeymd4   \nDiscord:                 tyksdfeymsbegxtzed460td5gw   \nHangouts Chat:           yksdfeymsbe6t0td5gzed4wgxt   \nMicrosoft Teams:         gwgxtzed4yksdfeymsbe6t0td5   \nMattermost:              wgxtzed4yksd5dfeymsbe6t0tg   \nWebhooks:                e6t0td5ykgwgxtzed4eymsbsdf   \nRegex Exclude:           ^0\\.1                        \nRegex Exclude Inverse:   ^0\\.3                        \nExclude Updated:         yes                          \n",
		},
		{
			name: "update all",
			args: []string{
				"--email", "weekly",
				"--slack", "mdsbe60td5gwgzetyksdfeyxt4",
				"--telegram", "sdfeyxt4mdsbe60td5gwgzetyk",
				"--discord", "4mdsbe60td5gwgzetyksdfeyxt",
				"--discord", "zext4mdsbe6tyksdfey0td5gwg",
				"--hangouts-chat", "etyksdfeyxt4mdsbe60td5gwgz",
				"--microsoft-teams", "0td5gwgzextbe6tyksdfey4mds",
				"--mattermost", "wgxtzed4yksd5dfeymsbe6t0tg",
				"--webhook", "tbe6tyksdfey4md0td5gwgzexs",
				"--regex-exclude", `^0\.1`,
				"--regex-exclude", `^0\.3-inverse`,
				"--exclude-prereleases",
				"--exclude-updated",
			},
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{minimalProject}),
			wantOutput:      "ID:                      mdsbe60td5gwgzetyksdfeyxt4                               \nName:                    golang/go                                                \nProvider:                github                                                   \nEmail:                   weekly                                                   \nSlack:                   mdsbe60td5gwgzetyksdfeyxt4                               \nTelegram:                sdfeyxt4mdsbe60td5gwgzetyk                               \nDiscord:                 4mdsbe60td5gwgzetyksdfeyxt, zext4mdsbe6tyksdfey0td5gwg   \nHangouts Chat:           etyksdfeyxt4mdsbe60td5gwgz                               \nMicrosoft Teams:         0td5gwgzextbe6tyksdfey4mds                               \nMattermost:              wgxtzed4yksd5dfeymsbe6t0tg                               \nWebhooks:                tbe6tyksdfey4md0td5gwgzexs                               \nRegex Exclude:           ^0\\.1                                                    \nRegex Exclude Inverse:   ^0\\.3                                                    \nExclude Pre-Releases:    yes                                                      \nExclude Updated:         yes                                                      \n",
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
					cmd.WithArgs(append([]string{"project", "update", "mdsbe60td5gwgzetyksdfeyxt4"}, tc.args...)...),
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
			t.Run("by name", func(t *testing.T) {
				var outputBuf bytes.Buffer
				if err := newCommand(t,
					cmd.WithArgs(append([]string{"project", "update", "github", "golang/go"}, tc.args...)...),
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
		})
	}
}
