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

func TestProjectCmd_Search(t *testing.T) {
	for _, tc := range []struct {
		name            string
		args            []string
		projectsService cmd.ProjectsService
		wantOutput      string
		wantError       error
	}{
		{
			name:            "no projects",
			args:            []string{"golang"},
			projectsService: newMockProjectsService(1, nil),
			wantOutput:      "No projects found.\n",
		},
		{
			name: "minimal project",
			args: []string{"golang"},
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{
				{ID: "mdsbe60td5gwgzetyksdfeyxt4", Name: "golang/go", Provider: "github"},
			}),
			wantOutput: "ID                           NAME        PROVIDER \nmdsbe60td5gwgzetyksdfeyxt4   golang/go   github     \n",
		},
		{
			name: "projects filter by provider",
			args: []string{"golang", "--provider", "github"},
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{
				{ID: "mdsbe60td5gwgzetyksdfeyxt4", Name: "golang/go", Provider: "github"},
				{ID: "ksdfeyxt4mdsbe60td5gwgzety", Name: "newreleases/cli-go", Provider: "github"},
				{ID: "gwgzetyksdfeyxt4mdsbe60td5", Name: "vue", Provider: "npm"},
			}),
			wantOutput: "ID                           NAME                 PROVIDER \nmdsbe60td5gwgzetyksdfeyxt4   golang/go            github     \nksdfeyxt4mdsbe60td5gwgzety   newreleases/cli-go   github     \n",
		},
		{
			name:            "full project",
			args:            []string{"golang"},
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{fullProject}),
			wantOutput:      "ID                           NAME        PROVIDER   EMAIL   SLACK                        TELEGRAM                     DISCORD                      HANGOUTS CHAT                MICROSOFT TEAMS              MATTERMOST                   ROCKET CHAT       MATRIX                  WEBHOOK                      REGEX EXCLUDE   REGEX EXCLUDE INVERSE   EXCLUDE PRE-RELEASES   EXCLUDE UPDATED   NOTE           TAGS         \nmdsbe60td5gwgzetyksdfeyxt4   golang/go   github     daily   zetyksdfeymdsbe60td5gwgxt4   sbe60td5gwgxtzetyksdfeymd4   tyksdfeymsbegxtzed460td5gw   yksdfeymsbe6t0td5gzed4wgxt   gwgxtzed4yksdfeymsbe6t0td5   wgxtzed4yksd5dfeymsbe6t0tg   5dfeymsbe6t0tgwgxtzed4yksd   4yksd5e6twgxtzdfeymsbed0tg   e6t0td5ykgwgxtzed4eymsbsdf   ^0\\.1           ^0\\.3                   yes                    yes               Initial no...   33f1db7254b9   \n",
		},
		{
			name:            "error",
			args:            []string{"golang"},
			projectsService: newMockProjectsService(1, errTest),
			wantError:       errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs(append([]string{"project", "search"}, tc.args...)...),
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
	}
}
