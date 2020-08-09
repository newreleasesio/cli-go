// Copyright (c) 2020, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file. {

package cmd_test

import (
	"context"
	"sort"

	"newreleases.io/newreleases"
)

var (
	minimalProject = newreleases.Project{
		ID:       "mdsbe60td5gwgzetyksdfeyxt4",
		Name:     "golang/go",
		Provider: "github",
	}
	fullProject = newreleases.Project{
		ID:                     "mdsbe60td5gwgzetyksdfeyxt4",
		Name:                   "golang/go",
		Provider:               "github",
		EmailNotification:      newreleases.EmailNotificationDaily,
		SlackIDs:               []string{"zetyksdfeymdsbe60td5gwgxt4"},
		TelegramChatIDs:        []string{"sbe60td5gwgxtzetyksdfeymd4"},
		DiscordIDs:             []string{"tyksdfeymsbegxtzed460td5gw"},
		HangoutsChatWebhookIDs: []string{"yksdfeymsbe6t0td5gzed4wgxt"},
		MSTeamsWebhookIDs:      []string{"gwgxtzed4yksdfeymsbe6t0td5"},
		MattermostWebhookIDs:   []string{"wgxtzed4yksd5dfeymsbe6t0tg"},
		WebhookIDs:             []string{"e6t0td5ykgwgxtzed4eymsbsdf"},
		Exclusions: []newreleases.Exclusion{
			{Value: `^0\.1`},
			{Value: `^0\.3`, Inverse: true},
		},
		ExcludePrereleases: true,
		ExcludeUpdated:     true,
	}
)

type mockProjectsService struct {
	pages    [][]newreleases.Project
	lastPage int
	err      error
}

func newMockProjectsService(lastPage int, err error, pages ...[]newreleases.Project) (s mockProjectsService) {
	return mockProjectsService{pages: pages, lastPage: lastPage, err: err}
}

func (s mockProjectsService) List(ctx context.Context, o newreleases.ProjectListOptions) (projects []newreleases.Project, lastPage int, err error) {
	if len(s.pages) == 0 {
		return nil, s.lastPage, s.err
	}
	if o.Provider != "" {
		for _, p := range s.pages[o.Page-1] {
			if p.Provider == o.Provider {
				projects = append(projects, p)
			}
		}
	} else {
		projects = s.pages[o.Page-1]
	}
	switch o.Order {
	case newreleases.ProjectListOrderName:
		sort.SliceStable(projects, func(i, j int) (less bool) {
			return projects[i].Name+projects[i].Provider < projects[j].Name+projects[j].Provider
		})
	case newreleases.ProjectListOrderAdded:
		sort.SliceStable(projects, func(i, j int) (less bool) {
			return projects[i].ID < projects[j].ID
		})
	}
	if o.Reverse {
		sort.Slice(projects, func(i, j int) (less bool) {
			return false
		})
	}
	return projects, s.lastPage, s.err
}

func (s mockProjectsService) Search(ctx context.Context, query, provider string) (projects []newreleases.Project, err error) {
	if len(s.pages) == 0 {
		return nil, s.err
	}
	if provider != "" {
		for _, p := range s.pages[0] {
			if p.Provider == provider {
				projects = append(projects, p)
			}
		}
	} else {
		projects = s.pages[0]
	}
	return projects, s.err
}

func (s mockProjectsService) GetByID(ctx context.Context, id string) (project *newreleases.Project, err error) {
	if len(s.pages) == 0 || len(s.pages[0]) == 0 {
		return nil, s.err
	}
	return &s.pages[0][0], s.err
}

func (s mockProjectsService) GetByName(ctx context.Context, provider, name string) (project *newreleases.Project, err error) {
	if len(s.pages) == 0 || len(s.pages[0]) == 0 {
		return nil, s.err
	}
	return &s.pages[0][0], s.err
}

func (s mockProjectsService) Add(ctx context.Context, provider, name string, o *newreleases.ProjectOptions) (project *newreleases.Project, err error) {
	project = &newreleases.Project{
		ID:       "new",
		Name:     name,
		Provider: provider,
		URL:      "url",
	}
	if o.EmailNotification != nil {
		project.EmailNotification = *o.EmailNotification
	}
	if o.SlackIDs != nil {
		project.SlackIDs = o.SlackIDs
	}
	if o.TelegramChatIDs != nil {
		project.TelegramChatIDs = o.TelegramChatIDs
	}
	if o.DiscordIDs != nil {
		project.DiscordIDs = o.DiscordIDs
	}
	if o.HangoutsChatWebhookIDs != nil {
		project.HangoutsChatWebhookIDs = o.HangoutsChatWebhookIDs
	}
	if o.MSTeamsWebhookIDs != nil {
		project.MSTeamsWebhookIDs = o.MSTeamsWebhookIDs
	}
	if o.MattermostWebhookIDs != nil {
		project.MattermostWebhookIDs = o.MattermostWebhookIDs
	}
	if o.WebhookIDs != nil {
		project.WebhookIDs = o.WebhookIDs
	}
	if o.Exclusions != nil {
		project.Exclusions = o.Exclusions
	}
	if o.ExcludePrereleases != nil {
		project.ExcludePrereleases = *o.ExcludePrereleases
	}
	if o.ExcludeUpdated != nil {
		project.ExcludeUpdated = *o.ExcludeUpdated
	}
	return project, s.err
}

func (s mockProjectsService) UpdateByID(ctx context.Context, id string, o *newreleases.ProjectOptions) (project *newreleases.Project, err error) {
	if len(s.pages) == 0 || len(s.pages[0]) == 0 {
		return nil, s.err
	}
	project = &s.pages[0][0]
	if o.EmailNotification != nil {
		project.EmailNotification = *o.EmailNotification
	}
	if o.SlackIDs != nil {
		project.SlackIDs = o.SlackIDs
	}
	if o.TelegramChatIDs != nil {
		project.TelegramChatIDs = o.TelegramChatIDs
	}
	if o.DiscordIDs != nil {
		project.DiscordIDs = o.DiscordIDs
	}
	if o.HangoutsChatWebhookIDs != nil {
		project.HangoutsChatWebhookIDs = o.HangoutsChatWebhookIDs
	}
	if o.MSTeamsWebhookIDs != nil {
		project.MSTeamsWebhookIDs = o.MSTeamsWebhookIDs
	}
	if o.MattermostWebhookIDs != nil {
		project.MattermostWebhookIDs = o.MattermostWebhookIDs
	}
	if o.WebhookIDs != nil {
		project.WebhookIDs = o.WebhookIDs
	}
	if o.Exclusions != nil {
		project.Exclusions = o.Exclusions
	}
	if o.ExcludePrereleases != nil {
		project.ExcludePrereleases = *o.ExcludePrereleases
	}
	if o.ExcludeUpdated != nil {
		project.ExcludeUpdated = *o.ExcludeUpdated
	}
	return project, s.err
}

func (s mockProjectsService) UpdateByName(ctx context.Context, provider, name string, o *newreleases.ProjectOptions) (project *newreleases.Project, err error) {
	if len(s.pages) == 0 || len(s.pages[0]) == 0 {
		return nil, s.err
	}
	project = &s.pages[0][0]
	if o.EmailNotification != nil {
		project.EmailNotification = *o.EmailNotification
	}
	if o.SlackIDs != nil {
		project.SlackIDs = o.SlackIDs
	}
	if o.TelegramChatIDs != nil {
		project.TelegramChatIDs = o.TelegramChatIDs
	}
	if o.DiscordIDs != nil {
		project.DiscordIDs = o.DiscordIDs
	}
	if o.HangoutsChatWebhookIDs != nil {
		project.HangoutsChatWebhookIDs = o.HangoutsChatWebhookIDs
	}
	if o.MSTeamsWebhookIDs != nil {
		project.MSTeamsWebhookIDs = o.MSTeamsWebhookIDs
	}
	if o.MattermostWebhookIDs != nil {
		project.MattermostWebhookIDs = o.MattermostWebhookIDs
	}
	if o.WebhookIDs != nil {
		project.WebhookIDs = o.WebhookIDs
	}
	if o.Exclusions != nil {
		project.Exclusions = o.Exclusions
	}
	if o.ExcludePrereleases != nil {
		project.ExcludePrereleases = *o.ExcludePrereleases
	}
	if o.ExcludeUpdated != nil {
		project.ExcludeUpdated = *o.ExcludeUpdated
	}
	return project, s.err
}

func (s mockProjectsService) DeleteByID(ctx context.Context, id string) (err error) {
	return s.err
}

func (s mockProjectsService) DeleteByName(ctx context.Context, provider, name string) (err error) {
	return s.err
}
