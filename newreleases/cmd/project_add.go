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

func (c *command) initProjectAddCmd(projectCmd *cobra.Command) (err error) {
	var (
		optionNameEmail              = "email"
		optionNameSlack              = "slack"
		optionNameTelegram           = "telegram"
		optionNameDiscord            = "discord"
		optionNameHangoutsChat       = "hangouts-chat"
		optionNameMicrosoftTeams     = "microsoft-teams"
		optionNameMattermost         = "mattermost"
		optionNameWebhook            = "webhook"
		optionNameExclusions         = "regex-exclude"
		optionNameExcludePrereleases = "exclude-prereleases"
		optionNameExcludeUpdated     = "exclude-updated"
	)

	cmd := &cobra.Command{
		Use:   "add PROVIDER PROJECT_NAME",
		Short: "Add a project to track",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			ctx, cancel := newClientContext(c.config)
			defer cancel()

			if len(args) != 2 {
				return cmd.Help()
			}

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
			o.SlackIDs, err = flags.GetStringArray(optionNameSlack)
			if err != nil {
				return err
			}
			o.TelegramChatIDs, err = flags.GetStringArray(optionNameTelegram)
			if err != nil {
				return err
			}
			o.DiscordIDs, err = flags.GetStringArray(optionNameDiscord)
			if err != nil {
				return err
			}
			o.HangoutsChatWebhookIDs, err = flags.GetStringArray(optionNameHangoutsChat)
			if err != nil {
				return err
			}
			o.MSTeamsWebhookIDs, err = flags.GetStringArray(optionNameMicrosoftTeams)
			if err != nil {
				return err
			}
			o.MattermostWebhookIDs, err = flags.GetStringArray(optionNameMattermost)
			if err != nil {
				return err
			}
			o.WebhookIDs, err = flags.GetStringArray(optionNameWebhook)
			if err != nil {
				return err
			}
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

			project, err := c.projectsService.Add(ctx, args[0], args[1], o)
			if err != nil {
				return err
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

	cmd.Flags().String(optionNameEmail, "none", "frequency of email notifications: instant, hourly, daily, weekly, none")
	cmd.Flags().StringArray(optionNameSlack, nil, "Slack channel ID")
	cmd.Flags().StringArray(optionNameTelegram, nil, "Telegram chat ID")
	cmd.Flags().StringArray(optionNameDiscord, nil, "Discord channel ID")
	cmd.Flags().StringArray(optionNameHangoutsChat, nil, "Hangouts Chat webhook ID")
	cmd.Flags().StringArray(optionNameMicrosoftTeams, nil, "Microsoft Teams webhook ID")
	cmd.Flags().StringArray(optionNameMattermost, nil, "Mattermost webhook ID")
	cmd.Flags().StringArray(optionNameWebhook, nil, "Webhook ID")
	cmd.Flags().StringArray(optionNameExclusions, nil, "Regex version exclusion, suffix with \"-inverse\" for inclusion")
	cmd.Flags().Bool(optionNameExcludePrereleases, false, "exclude pre-releases")
	cmd.Flags().Bool(optionNameExcludeUpdated, false, "exclude updated")

	projectCmd.AddCommand(cmd)
	return addClientFlags(cmd)
}
