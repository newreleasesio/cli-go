// Copyright (c) 2020, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"strings"

	"github.com/spf13/cobra"
	"newreleases.io/newreleases"
)

func (c *command) initProjectUpdateCmd(projectCmd *cobra.Command) (err error) {
	var (
		optionNameEmail                = "email"
		optionNameSlack                = "slack"
		optionNameSlackRemove          = "slack-remove"
		optionNameTelegram             = "telegram"
		optionNameTelegramRemove       = "telegram-remove"
		optionNameDiscord              = "discord"
		optionNameDiscordRemove        = "discord-remove"
		optionNameHangoutsChat         = "hangouts-chat"
		optionNameHangoutsChatRemove   = "hangouts-chat-remove"
		optionNameMicrosoftTeams       = "microsoft-teams"
		optionNameMicrosoftTeamsRemove = "microsoft-teams-remove"
		optionNameWebhook              = "webhook"
		optionNameWebhookRemove        = "webhook-remove"
		optionNameExclusions           = "regex-exclude"
		optionNameExclusionsRemove     = "regex-exclude-remove"
		optionNameExcludePrereleases   = "exclude-prereleases"
		optionNameExcludeUpdated       = "exclude-updated"
	)

	cmd := &cobra.Command{
		Use:   "update [PROVIDER PROJECT_NAME] | [PROJECT_ID]",
		Short: "Update a tracked project",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			o := &newreleases.ProjectOptions{}

			flags := cmd.Flags()
			email, err := flags.GetString(optionNameEmail)
			if err != nil {
				return err
			}
			if email != "" {
				e := newreleases.EmailNotification(email)
				o.EmailNotification = &e
			}
			slackRemove, err := flags.GetBool(optionNameSlackRemove)
			if err != nil {
				return err
			}
			if slackRemove {
				o.SlackIDs = make([]string, 0)
			} else {
				slackIDs, err := flags.GetStringArray(optionNameSlack)
				if err != nil {
					return err
				}
				if len(slackIDs) > 0 {
					o.SlackIDs = slackIDs
				}
			}
			telegramRemove, err := flags.GetBool(optionNameTelegramRemove)
			if err != nil {
				return err
			}
			if telegramRemove {
				o.TelegramChatIDs = make([]string, 0)
			} else {
				telegramChatIDs, err := flags.GetStringArray(optionNameTelegram)
				if err != nil {
					return err
				}
				if len(telegramChatIDs) > 0 {
					o.TelegramChatIDs = telegramChatIDs
				}
			}
			discordRemove, err := flags.GetBool(optionNameDiscordRemove)
			if err != nil {
				return err
			}
			if discordRemove {
				o.DiscordIDs = make([]string, 0)
			} else {
				discordIDs, err := flags.GetStringArray(optionNameDiscord)
				if err != nil {
					return err
				}
				if len(discordIDs) > 0 {
					o.DiscordIDs = discordIDs
				}
			}
			hangoutsChatRemove, err := flags.GetBool(optionNameHangoutsChatRemove)
			if err != nil {
				return err
			}
			if hangoutsChatRemove {
				o.HangoutsChatWebhookIDs = make([]string, 0)
			} else {
				hangoutsChatWebhookIDs, err := flags.GetStringArray(optionNameHangoutsChat)
				if err != nil {
					return err
				}
				if len(hangoutsChatWebhookIDs) > 0 {
					o.HangoutsChatWebhookIDs = hangoutsChatWebhookIDs
				}
			}
			microsoftTeamsRemove, err := flags.GetBool(optionNameMicrosoftTeamsRemove)
			if err != nil {
				return err
			}
			if microsoftTeamsRemove {
				o.MSTeamsWebhookIDs = make([]string, 0)
			} else {
				msTeamsWebhookIDs, err := flags.GetStringArray(optionNameMicrosoftTeams)
				if err != nil {
					return err
				}
				if len(msTeamsWebhookIDs) > 0 {
					o.MSTeamsWebhookIDs = msTeamsWebhookIDs
				}
			}
			webhookRemove, err := flags.GetBool(optionNameWebhookRemove)
			if err != nil {
				return err
			}
			if webhookRemove {
				o.WebhookIDs = make([]string, 0)
			} else {
				webhookIDs, err := flags.GetStringArray(optionNameWebhook)
				if err != nil {
					return err
				}
				if len(webhookIDs) > 0 {
					o.WebhookIDs = webhookIDs
				}
			}
			exclusionsRemove, err := flags.GetBool(optionNameExclusionsRemove)
			if err != nil {
				return err
			}
			if exclusionsRemove {
				o.Exclusions = make([]newreleases.Exclusion, 0)
			} else {
				exclusions, err := flags.GetStringArray(optionNameExclusions)
				if err != nil {
					return err
				}
				for _, v := range exclusions {
					var inverse bool
					if strings.HasSuffix(v, "-inverse") {
						inverse = true
						v = strings.TrimSuffix(v, "-inverse")
					}
					o.Exclusions = append(o.Exclusions, newreleases.Exclusion{
						Value:   v,
						Inverse: inverse,
					})
				}
			}
			if flags.Changed(optionNameExcludePrereleases) {
				excludePrereleases, err := flags.GetBool(optionNameExcludePrereleases)
				if err != nil {
					return err
				}
				o.ExcludePrereleases = &excludePrereleases
			}
			if flags.Changed(optionNameExcludeUpdated) {
				excludeUpdated, err := flags.GetBool(optionNameExcludeUpdated)
				if err != nil {
					return err
				}
				o.ExcludeUpdated = &excludeUpdated
			}

			var project *newreleases.Project
			switch len(args) {
			case 1:
				project, err = c.projectsService.UpdateByID(ctx, args[0], o)
			case 2:
				project, err = c.projectsService.UpdateByName(ctx, args[0], args[1], o)
			default:
				return cmd.Help()
			}
			if err != nil {
				return err
			}

			if project == nil || err == newreleases.ErrNotFound {
				cmd.Println("Project not found.")
				return nil
			}

			printProject(cmd, project)
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if err := addClientConfigOptions(cmd, c.config); err != nil {
				return err
			}
			return c.setProjectsService(cmd, args)
		},
	}

	cmd.Flags().String(optionNameEmail, "none", "frequency of email notifications: hourly, daily, weekly, none")
	cmd.Flags().StringArray(optionNameSlack, nil, "Slack channel ID")
	cmd.Flags().Bool(optionNameSlackRemove, false, "remove Slack notifications")
	cmd.Flags().StringArray(optionNameTelegram, nil, "Telegram chat ID")
	cmd.Flags().Bool(optionNameTelegramRemove, false, "remove Telegram notifications")
	cmd.Flags().StringArray(optionNameDiscord, nil, "Discord channel ID")
	cmd.Flags().Bool(optionNameDiscordRemove, false, "remove Discord notifications")
	cmd.Flags().StringArray(optionNameHangoutsChat, nil, "Hangouts Chat webhook ID")
	cmd.Flags().Bool(optionNameHangoutsChatRemove, false, "remove Hangouts Chat notifications")
	cmd.Flags().StringArray(optionNameMicrosoftTeams, nil, "Microsoft Teams webhook ID")
	cmd.Flags().Bool(optionNameMicrosoftTeamsRemove, false, "remove Microsoft Teams notifications")
	cmd.Flags().StringArray(optionNameWebhook, nil, "Webhook ID")
	cmd.Flags().Bool(optionNameWebhookRemove, false, "remove Webhook notifications")
	cmd.Flags().StringArray(optionNameExclusions, nil, "Regex version exclusion, suffix with \"-inverse\" for inclusion")
	cmd.Flags().Bool(optionNameExclusionsRemove, false, "remove Regex version exclusions")
	cmd.Flags().Bool(optionNameExcludePrereleases, false, "exclude pre-releases")
	cmd.Flags().Bool(optionNameExcludeUpdated, false, "exclude updated")

	projectCmd.AddCommand(cmd)
	return addClientFlags(cmd)
}
