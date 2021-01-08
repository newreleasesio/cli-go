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

func TestProjectCmd_List(t *testing.T) {
	for _, tc := range []struct {
		name            string
		args            []string
		projectsService cmd.ProjectsService
		wantOutput      string
		wantError       error
	}{
		{
			name:            "no projects",
			projectsService: newMockProjectsService(1, nil),
			wantOutput:      "No projects found.\n",
		},
		{
			name: "minimal project",
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{
				{ID: "mdsbe60td5gwgzetyksdfeyxt4", Name: "golang/go", Provider: "github"},
			}),
			wantOutput: "ID                           NAME        PROVIDER \nmdsbe60td5gwgzetyksdfeyxt4   golang/go   github     \n",
		},
		{
			name: "projects page 2",
			args: []string{"--page", "2"},
			projectsService: newMockProjectsService(1, nil,
				[]newreleases.Project{
					{ID: "mdsbe60td5gwgzetyksdfeyxt4", Name: "golang/go", Provider: "github"},
				},
				[]newreleases.Project{
					{ID: "myxtdsbe60td5gwgzetyksdfe4", Name: "newreleases/cli-go", Provider: "github"},
				},
			),
			wantOutput: "ID                           NAME                 PROVIDER \nmyxtdsbe60td5gwgzetyksdfe4   newreleases/cli-go   github     \n",
		},
		{
			name: "projects filter by provider",
			args: []string{"--provider", "github"},
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{
				{ID: "mdsbe60td5gwgzetyksdfeyxt4", Name: "golang/go", Provider: "github"},
				{ID: "ksdfeyxt4mdsbe60td5gwgzety", Name: "newreleases/cli-go", Provider: "github"},
				{ID: "gwgzetyksdfeyxt4mdsbe60td5", Name: "vue", Provider: "npm"},
			}),
			wantOutput: "ID                           NAME                 PROVIDER \nmdsbe60td5gwgzetyksdfeyxt4   golang/go            github     \nksdfeyxt4mdsbe60td5gwgzety   newreleases/cli-go   github     \n",
		},
		{
			name: "projects order by name",
			args: []string{"--order", "name"},
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{
				{ID: "gwgzetyksdfeyxt4mdsbe60td5", Name: "vue", Provider: "npm"},
				{ID: "mdsbe60td5gwgzetyksdfeyxt4", Name: "golang/go", Provider: "github"},
				{ID: "ksdfeyxt4mdsbe60td5gwgzety", Name: "newreleases/cli-go", Provider: "github"},
			}),
			wantOutput: "ID                           NAME                 PROVIDER \nmdsbe60td5gwgzetyksdfeyxt4   golang/go            github     \nksdfeyxt4mdsbe60td5gwgzety   newreleases/cli-go   github     \ngwgzetyksdfeyxt4mdsbe60td5   vue                  npm        \n",
		},
		{
			name: "projects order by added time",
			args: []string{"--order", "added"},
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{
				{ID: "gwgzetyksdfeyxt4mdsbe60td5", Name: "vue", Provider: "npm"},
				{ID: "mdsbe60td5gwgzetyksdfeyxt4", Name: "golang/go", Provider: "github"},
				{ID: "ksdfeyxt4mdsbe60td5gwgzety", Name: "newreleases/cli-go", Provider: "github"},
			}),
			wantOutput: "ID                           NAME                 PROVIDER \ngwgzetyksdfeyxt4mdsbe60td5   vue                  npm        \nksdfeyxt4mdsbe60td5gwgzety   newreleases/cli-go   github     \nmdsbe60td5gwgzetyksdfeyxt4   golang/go            github     \n",
		},
		{
			name:            "full project",
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{fullProject}),
			wantOutput:      "ID                           NAME        PROVIDER   EMAIL   SLACK                        TELEGRAM                     DISCORD                      HANGOUTS CHAT                MICROSOFT TEAMS              MATTERMOST                   ROCKET CHAT                  WEBHOOK                      REGEX EXCLUDE   REGEX EXCLUDE INVERSE   EXCLUDE PRE-RELEASES   EXCLUDE UPDATED \nmdsbe60td5gwgzetyksdfeyxt4   golang/go   github     daily   zetyksdfeymdsbe60td5gwgxt4   sbe60td5gwgxtzetyksdfeymd4   tyksdfeymsbegxtzed460td5gw   yksdfeymsbe6t0td5gzed4wgxt   gwgxtzed4yksdfeymsbe6t0td5   wgxtzed4yksd5dfeymsbe6t0tg   5dfeymsbe6t0tgwgxtzed4yksd   e6t0td5ykgwgxtzed4eymsbsdf   ^0\\.1           ^0\\.3                   yes                    yes               \n",
		},
		{
			name:            "error",
			projectsService: newMockProjectsService(1, errTest),
			wantError:       errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs(append([]string{"project", "list"}, tc.args...)...),
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
