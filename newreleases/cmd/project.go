// Copyright (c) 2020, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"strings"

	"github.com/spf13/cobra"
	"newreleases.io/newreleases"
)

func (c *command) initProjectCmd() (err error) {
	cmd := &cobra.Command{
		Use:   "project",
		Short: "Manage tracked projects",
	}

	if err := c.initProjectListCmd(cmd); err != nil {
		return err
	}
	if err := c.initProjectSearchCmd(cmd); err != nil {
		return err
	}
	if err := c.initProjectGetCmd(cmd); err != nil {
		return err
	}
	if err := c.initProjectAddCmd(cmd); err != nil {
		return err
	}
	if err := c.initProjectUpdateCmd(cmd); err != nil {
		return err
	}
	if err := c.initProjectRemoveCmd(cmd); err != nil {
		return err
	}

	c.root.AddCommand(cmd)
	return nil
}

func (c *command) setProjectsService(cmd *cobra.Command, args []string) (err error) {
	if c.projectsService != nil {
		return nil
	}
	client, err := c.getClient(cmd)
	if err != nil {
		return err
	}
	c.projectsService = client.Projects
	return nil
}

type projectsService interface {
	List(ctx context.Context, o newreleases.ProjectListOptions) (projects []newreleases.Project, lastPage int, err error)
	Search(ctx context.Context, query, provider string) (projects []newreleases.Project, err error)
	GetByID(ctx context.Context, id string) (project *newreleases.Project, err error)
	GetByName(ctx context.Context, provider, name string) (project *newreleases.Project, err error)
	Add(ctx context.Context, provider, name string, o *newreleases.ProjectOptions) (project *newreleases.Project, err error)
	UpdateByID(ctx context.Context, id string, o *newreleases.ProjectOptions) (project *newreleases.Project, err error)
	UpdateByName(ctx context.Context, provider, name string, o *newreleases.ProjectOptions) (project *newreleases.Project, err error)
	DeleteByID(ctx context.Context, id string) (err error)
	DeleteByName(ctx context.Context, provider, name string) (err error)
}

func printProjectsTable(cmd *cobra.Command, projects []newreleases.Project) {
	table := newTable(cmd.OutOrStdout())

	var (
		hasEmailNotification  bool
		hasSlack              bool
		hasTelegram           bool
		hasDiscord            bool
		hasHangoutsChat       bool
		hasMicrosoftTeams     bool
		hasMattermost         bool
		hasWebhook            bool
		hasExclusions         bool
		hasInclusions         bool
		hasExcludePrereleases bool
		hasExcludeUpdated     bool
	)
	for _, p := range projects {
		if p.EmailNotification != newreleases.EmailNotificationNone && p.EmailNotification != "" {
			hasEmailNotification = true
		}
		if len(p.SlackIDs) > 0 {
			hasSlack = true
		}
		if len(p.TelegramChatIDs) > 0 {
			hasTelegram = true
		}
		if len(p.DiscordIDs) > 0 {
			hasDiscord = true
		}
		if len(p.HangoutsChatWebhookIDs) > 0 {
			hasHangoutsChat = true
		}
		if len(p.MSTeamsWebhookIDs) > 0 {
			hasMicrosoftTeams = true
		}
		if len(p.MattermostWebhookIDs) > 0 {
			hasMattermost = true
		}
		if len(p.WebhookIDs) > 0 {
			hasWebhook = true
		}
		if len(p.Exclusions) > 0 {
			for _, e := range p.Exclusions {
				if e.Inverse {
					hasInclusions = true
				} else {
					hasExclusions = true
				}
			}
		}
		if p.ExcludePrereleases {
			hasExcludePrereleases = true
		}
		if p.ExcludeUpdated {
			hasExcludeUpdated = true
		}
	}

	header := []string{
		"ID",
		"Name",
		"Provider",
	}
	if hasEmailNotification {
		header = append(header, "Email")
	}
	if hasSlack {
		header = append(header, "Slack")
	}
	if hasTelegram {
		header = append(header, "Telegram")
	}
	if hasDiscord {
		header = append(header, "Discord")
	}
	if hasHangoutsChat {
		header = append(header, "Hangouts Chat")
	}
	if hasMicrosoftTeams {
		header = append(header, "Microsoft Teams")
	}
	if hasMattermost {
		header = append(header, "Mattermost")
	}
	if hasWebhook {
		header = append(header, "Webhook")
	}
	if hasExclusions {
		header = append(header, "Regex Exclude")
	}
	if hasInclusions {
		header = append(header, "Regex Exclude Inverse")
	}
	if hasExcludePrereleases {
		header = append(header, "Exclude Pre-Releases")
	}
	if hasExcludeUpdated {
		header = append(header, "Exclude Updated")
	}
	table.SetHeader(header)
	for _, p := range projects {
		r := []string{p.ID, p.Name, p.Provider}
		if hasEmailNotification {
			r = append(r, string(p.EmailNotification))
		}
		if hasSlack {
			r = append(r, strings.Join(p.SlackIDs, ", "))
		}
		if hasTelegram {
			r = append(r, strings.Join(p.TelegramChatIDs, ", "))
		}
		if hasDiscord {
			r = append(r, strings.Join(p.DiscordIDs, ", "))
		}
		if hasHangoutsChat {
			r = append(r, strings.Join(p.HangoutsChatWebhookIDs, ", "))
		}
		if hasMicrosoftTeams {
			r = append(r, strings.Join(p.MSTeamsWebhookIDs, ", "))
		}
		if hasMattermost {
			r = append(r, strings.Join(p.MattermostWebhookIDs, ", "))
		}
		if hasWebhook {
			r = append(r, strings.Join(p.WebhookIDs, ", "))
		}
		if hasExclusions {
			var l []string
			for _, e := range p.Exclusions {
				if !e.Inverse {
					l = append(l, e.Value)
				}
			}
			r = append(r, strings.Join(l, ", "))
		}
		if hasInclusions {
			var l []string
			for _, e := range p.Exclusions {
				if e.Inverse {
					l = append(l, e.Value)
				}
			}
			r = append(r, strings.Join(l, ", "))
		}
		if hasExcludePrereleases {
			r = append(r, yesNo(p.ExcludePrereleases))
		}
		if hasExcludeUpdated {
			r = append(r, yesNo(p.ExcludeUpdated))
		}
		table.Append(r)
	}
	table.Render()
}

func printProject(cmd *cobra.Command, p *newreleases.Project) {
	table := newTable(cmd.OutOrStdout())
	table.Append([]string{"ID:", p.ID})
	table.Append([]string{"Name:", p.Name})
	table.Append([]string{"Provider:", p.Provider})
	if p.EmailNotification != newreleases.EmailNotificationNone && p.EmailNotification != "" {
		table.Append([]string{"Email:", string(p.EmailNotification)})
	}
	if len(p.SlackIDs) > 0 {
		table.Append([]string{"Slack:", strings.Join(p.SlackIDs, ", ")})
	}
	if len(p.TelegramChatIDs) > 0 {
		table.Append([]string{"Telegram:", strings.Join(p.TelegramChatIDs, ", ")})
	}
	if len(p.DiscordIDs) > 0 {
		table.Append([]string{"Discord:", strings.Join(p.DiscordIDs, ", ")})
	}
	if len(p.HangoutsChatWebhookIDs) > 0 {
		table.Append([]string{"Hangouts Chat:", strings.Join(p.HangoutsChatWebhookIDs, ", ")})
	}
	if len(p.MSTeamsWebhookIDs) > 0 {
		table.Append([]string{"Microsoft Teams:", strings.Join(p.MSTeamsWebhookIDs, ", ")})
	}
	if len(p.MattermostWebhookIDs) > 0 {
		table.Append([]string{"Mattermost:", strings.Join(p.MattermostWebhookIDs, ", ")})
	}
	if len(p.WebhookIDs) > 0 {
		table.Append([]string{"Webhooks:", strings.Join(p.WebhookIDs, ", ")})
	}
	var excluded, excludedInverse []string
	for _, e := range p.Exclusions {
		if e.Inverse {
			excludedInverse = append(excludedInverse, e.Value)
		} else {
			excluded = append(excluded, e.Value)
		}
	}
	if len(excluded) > 0 {
		table.Append([]string{"Regex Exclude:", strings.Join(excluded, ", ")})
	}
	if len(excludedInverse) > 0 {
		table.Append([]string{"Regex Exclude Inverse:", strings.Join(excludedInverse, ", ")})
	}
	if p.ExcludePrereleases {
		table.Append([]string{"Exclude Pre-Releases:", "yes"})
	}
	if p.ExcludeUpdated {
		table.Append([]string{"Exclude Updated:", "yes"})
	}
	table.Render()
}
