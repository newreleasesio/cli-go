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

func TestProjectCmd_Remove(t *testing.T) {
	for _, tc := range []struct {
		name            string
		projectsService cmd.ProjectsService
		wantOutput      string
		wantError       error
	}{
		{
			name:            "no project",
			projectsService: newMockProjectsService(1, newreleases.ErrNotFound),
			wantOutput:      "Project not found.\n",
		},
		{
			name: "project",
			projectsService: newMockProjectsService(1, nil, []newreleases.Project{minimalProject}),
			wantOutput: "",
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
					cmd.WithArgs("project", "remove", "mdsbe60td5gwgzetyksdfeyxt4"),
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
					cmd.WithArgs("project", "remove", "github", "golang/go"),
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
