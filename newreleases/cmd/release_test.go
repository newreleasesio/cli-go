// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd_test

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"newreleases.io/cmd/newreleases/cmd"
	"newreleases.io/newreleases"
)

func TestReleaseCmd_List(t *testing.T) {
	for _, tc := range []struct {
		name            string
		args            []string
		releasesService cmd.ReleasesService
		wantOutputFunc  func() string
		wantError       error
	}{
		{
			name:            "no releases",
			args:            []string{"github", "golang/go"},
			releasesService: newMockReleasesService(nil, nil, 1, nil),
			wantOutputFunc:  func() string { return "No releases found.\n" },
		},
		{
			name:            "no releases page 2",
			args:            []string{"github", "golang/go", "--page", "2"},
			releasesService: newMockReleasesService(nil, nil, 1, nil),
			wantOutputFunc:  func() string { return "No releases found on page 2.\n" },
		},
		{
			name:            "no releases page 2",
			args:            []string{"github", "golang/go", "--page", "2"},
			releasesService: newMockReleasesService(nil, nil, 1, nil),
			wantOutputFunc:  func() string { return "No releases found on page 2.\n" },
		},
		{
			name: "with releases",
			args: []string{"github", "golang/go"},
			releasesService: newMockReleasesService([]newreleases.Release{
				{Version: "v1.25.0", Date: newTime(t, "2019-10-22T01:45:55Z")},
				{Version: "v1.21.6", Date: newTime(t, "2019-09-21T11:25:00Z"), IsPrerelease: true},
				{Version: "v1.21.5", Date: newTime(t, "2019-09-20T01:03:00Z"), HasNote: true},
				{Version: "v1.20.0", Date: newTime(t, "2019-09-01T15:14:00Z"), IsUpdated: true},
				{Version: "v1.18.88", Date: newTime(t, "2019-08-11T19:57:01Z"), IsExcluded: true},
			}, nil, 1, nil),
			wantOutputFunc: func() string {
				return fmt.Sprintf("VERSION    DATE                             PRE-RELEASE   HAS NOTE   UPDATED   EXCLUDED \nv1.25.0    %s   no            no         no        no         \nv1.21.6    %s   yes           no         no        no         \nv1.21.5    %s   no            yes        no        no         \nv1.20.0    %s   no            no         yes       no         \nv1.18.88   %s   no            no         no        yes        \n",
					newTime(t, "2019-10-22T01:45:55Z").Local(),
					newTime(t, "2019-09-21T11:25:00Z").Local(),
					newTime(t, "2019-09-20T01:03:00Z").Local(),
					newTime(t, "2019-09-01T15:14:00Z").Local(),
					newTime(t, "2019-08-11T19:57:01Z").Local(),
				)
			},
		},
		{
			name:            "error",
			args:            []string{"github", "golang/go"},
			releasesService: newMockReleasesService(nil, nil, 1, errTest),
			wantOutputFunc:  func() string { return "" },
			wantError:       errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs(append([]string{"release", "list"}, tc.args...)...),
				cmd.WithOutput(&outputBuf),
				cmd.WithReleasesService(tc.releasesService),
			).Execute(); err != tc.wantError {
				t.Fatalf("got error %v, want %v", err, tc.wantError)
			}

			gotOutput := outputBuf.String()
			if wantOutput := tc.wantOutputFunc(); gotOutput != wantOutput {
				t.Errorf("got output %q, want %q", gotOutput, wantOutput)
			}
		})
	}
}

func TestReleaseCmd_Get(t *testing.T) {
	for _, tc := range []struct {
		name            string
		args            []string
		releasesService cmd.ReleasesService
		wantOutputFunc  func() string
		wantError       error
	}{
		{
			name:            "no releases",
			args:            []string{"github", "golang/go", "v0.1.0"},
			releasesService: newMockReleasesService(nil, nil, 1, nil),
			wantOutputFunc:  func() string { return "Release not found.\n" },
		},
		{
			name: "release",
			args: []string{"github", "golang/go", "v0.1.0"},
			releasesService: newMockReleasesService([]newreleases.Release{
				{Version: "v0.1.0", Date: newTime(t, "2019-10-22T01:45:55Z")},
			}, nil, 1, nil),
			wantOutputFunc: func() string {
				return fmt.Sprintf("Version:       v0.1.0                           \nDate:          %s   \nPre-Release:   no                               \nHas Note:      no                               \nUpdated:       no                               \nExcluded:      no                               \n",
					newTime(t, "2019-10-22T01:45:55Z").Local(),
				)
			},
		},
		{
			name:            "error",
			args:            []string{"github", "golang/go", "v0.1.0"},
			releasesService: newMockReleasesService(nil, nil, 1, errTest),
			wantOutputFunc:  func() string { return "" },
			wantError:       errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs(append([]string{"release", "get"}, tc.args...)...),
				cmd.WithOutput(&outputBuf),
				cmd.WithReleasesService(tc.releasesService),
			).Execute(); err != tc.wantError {
				t.Fatalf("got error %v, want %v", err, tc.wantError)
			}

			gotOutput := outputBuf.String()
			if wantOutput := tc.wantOutputFunc(); gotOutput != wantOutput {
				t.Errorf("got output %q, want %q", gotOutput, wantOutput)
			}
		})
	}
}

func TestReleaseCmd_Note(t *testing.T) {
	for _, tc := range []struct {
		name            string
		args            []string
		releasesService cmd.ReleasesService
		wantOutput      string
		wantError       error
	}{
		{
			name:            "no notes",
			args:            []string{"github", "golang/go", "v0.1.0"},
			releasesService: newMockReleasesService(nil, nil, 1, nil),
			wantOutput:      "Release note not found.\n",
		},
		{
			name: "note",
			args: []string{"github", "golang/go", "v0.1.0"},
			releasesService: newMockReleasesService(nil, &newreleases.ReleaseNote{
				Title:   "Some awesome new release",
				Message: "<ul><li>Everything just works</li><li>What else to say?</li></ul>",
				URL:     "https://github.com/newreleases/cli-go/releases",
			}, 1, nil),
			wantOutput: "Some awesome new release\n\n* Everything just works\n* What else to say?\n\nhttps://github.com/newreleases/cli-go/releases\n\n",
		},
		{
			name:            "error",
			args:            []string{"github", "golang/go", "v0.1.0"},
			releasesService: newMockReleasesService(nil, nil, 1, errTest),
			wantError:       errTest,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var outputBuf bytes.Buffer
			if err := newCommand(t,
				cmd.WithArgs(append([]string{"release", "note"}, tc.args...)...),
				cmd.WithOutput(&outputBuf),
				cmd.WithReleasesService(tc.releasesService),
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

type mockReleasesService struct {
	releases []newreleases.Release
	note     *newreleases.ReleaseNote
	lastPage int
	err      error
}

func newMockReleasesService(releases []newreleases.Release, note *newreleases.ReleaseNote, lastPage int, err error) (s mockReleasesService) {
	return mockReleasesService{releases: releases, note: note, lastPage: lastPage, err: err}
}

func (s mockReleasesService) ListByProjectID(ctx context.Context, projectID string, page int) (releases []newreleases.Release, lastPage int, err error) {
	return s.releases, s.lastPage, s.err
}

func (s mockReleasesService) ListByProjectName(ctx context.Context, provider, projectName string, page int) (releases []newreleases.Release, lastPage int, err error) {
	return s.releases, s.lastPage, s.err
}

func (s mockReleasesService) GetByProjectID(ctx context.Context, projectID, version string) (release *newreleases.Release, err error) {
	if len(s.releases) == 0 {
		return nil, s.err
	}
	return &s.releases[0], s.err
}

func (s mockReleasesService) GetByProjectName(ctx context.Context, provider, projectName, version string) (release *newreleases.Release, err error) {
	if len(s.releases) == 0 {
		return nil, s.err
	}
	return &s.releases[0], s.err
}

func (s mockReleasesService) GetNoteByProjectID(ctx context.Context, projectID string, version string) (release *newreleases.ReleaseNote, err error) {
	return s.note, s.err
}

func (s mockReleasesService) GetNoteByProjectName(ctx context.Context, provider, projectName string, version string) (release *newreleases.ReleaseNote, err error) {
	return s.note, s.err
}
