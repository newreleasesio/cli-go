// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import "io"

type (
	Command                       = command
	Option                        = option
	PasswordReader                = passwordReader
	AuthKeysGetter                = authKeysGetter
	AuthService                   = authService
	ProvidersService              = providersService
	ProjectsService               = projectsService
	ReleasesService               = releasesService
	SlackChannelsService          = slackChannelsService
	TelegramChatsService          = telegramChatsService
	DiscordChannelsService        = discordChannelsService
	HangoutsChatWebhooksService   = hangoutsChatWebhooksService
	MicrosoftTeamsWebhooksService = microsoftTeamsWebhooksService
	MattermostWebhooksService     = mattermostWebhooksService
	RocketchatWebhooksService     = rocketchatWebhooksService
	MatrixRoomsService            = matrixRoomsService
	WebhooksService               = webhooksService
	TagsService                   = tagsService
)

var (
	NewCommand = newCommand
)

func WithCfgFile(f string) func(c *Command) {
	return func(c *Command) {
		c.cfgFile = f
	}
}

func WithHomeDir(dir string) func(c *Command) {
	return func(c *Command) {
		c.homeDir = dir
	}
}

func WithArgs(a ...string) func(c *Command) {
	return func(c *Command) {
		c.root.SetArgs(a)
	}
}

func WithInput(r io.Reader) func(c *Command) {
	return func(c *Command) {
		c.root.SetIn(r)
	}
}

func WithOutput(w io.Writer) func(c *Command) {
	return func(c *Command) {
		c.root.SetOut(w)
	}
}

func WithErrorOutput(w io.Writer) func(c *Command) {
	return func(c *Command) {
		c.root.SetErr(w)
	}
}

func WithPasswordReader(r PasswordReader) func(c *Command) {
	return func(c *Command) {
		c.passwordReader = r
	}
}

func WithAuthKeysGetter(g AuthKeysGetter) func(c *Command) {
	return func(c *Command) {
		c.authKeysGetter = g
	}
}

func WithAuthService(s AuthService) func(c *Command) {
	return func(c *Command) {
		c.authService = s
	}
}

func WithProvidersService(s ProvidersService) func(c *Command) {
	return func(c *Command) {
		c.providersService = s
	}
}

func WithProjectsService(s ProjectsService) func(c *Command) {
	return func(c *Command) {
		c.projectsService = s
	}
}

func WithReleasesService(s ReleasesService) func(c *Command) {
	return func(c *Command) {
		c.releasesService = s
	}
}

func WithSlackChannelsService(s SlackChannelsService) func(c *Command) {
	return func(c *Command) {
		c.slackChannelsService = s
	}
}

func WithTelegramChatsService(s TelegramChatsService) func(c *Command) {
	return func(c *Command) {
		c.telegramChatsService = s
	}
}

func WithDiscordChannelsService(s DiscordChannelsService) func(c *Command) {
	return func(c *Command) {
		c.discordChannelsService = s
	}
}

func WithHangoutsChatWebhooksService(s HangoutsChatWebhooksService) func(c *Command) {
	return func(c *Command) {
		c.hangoutsChatWebhooksService = s
	}
}

func WithMicrosoftTeamsWebhooksService(s MicrosoftTeamsWebhooksService) func(c *Command) {
	return func(c *Command) {
		c.microsoftTeamsWebhooksService = s
	}
}

func WithMattermostWebhooksService(s MattermostWebhooksService) func(c *Command) {
	return func(c *Command) {
		c.mattermostWebhooksService = s
	}
}

func WithRocketchatWebhooksService(s RocketchatWebhooksService) func(c *Command) {
	return func(c *Command) {
		c.rocketchatWebhooksService = s
	}
}

func WithMatrixRoomsService(s MatrixRoomsService) func(c *Command) {
	return func(c *Command) {
		c.matrixRoomsService = s
	}
}

func WithWebhooksService(s WebhooksService) func(c *Command) {
	return func(c *Command) {
		c.webhooksService = s
	}
}

func WithTagsService(s TagsService) func(c *Command) {
	return func(c *Command) {
		c.tagsService = s
	}
}
