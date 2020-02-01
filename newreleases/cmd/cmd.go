// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"newreleases.io/newreleases"
)

const (
	optionNameAuthKey     = "auth-key"
	optionNameTimeout     = "timeout"
	optionNameAPIEndpoint = "api-endpoint"
)

func init() {
	cobra.EnableCommandSorting = false
}

type command struct {
	root                          *cobra.Command
	config                        *viper.Viper
	client                        *newreleases.Client
	cfgFile                       string
	homeDir                       string
	passwordReader                passwordReader
	authKeysGetter                authKeysGetter
	authService                   authService
	projectsService               projectsService
	releasesService               releasesService
	providersService              providersService
	slackChannelsService          slackChannelsService
	telegramChatsService          telegramChatsService
	discordChannelsService        discordChannelsService
	hangoutsChatWebhooksService   hangoutsChatWebhooksService
	microsoftTeamsWebhooksService microsoftTeamsWebhooksService
	webhooksService               webhooksService
}

type option func(*command)

func newCommand(opts ...option) (c *command, err error) {
	c = &command{
		root: &cobra.Command{
			Use:   "newreleases",
			Short: "newreleases manages projects on NewReleases service",
			Long: `NewReleases is a release tracker for software engineers.

More information at https://newreleases.io.`,
			SilenceErrors: true,
			SilenceUsage:  true,
			PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
				cmdName := cmd.Name()
				return c.initConfig(cmdName != cmdNameConfigure && cmdName != cmdNameGetAuthKey)
			},
		},
	}

	for _, o := range opts {
		o(c)
	}
	if c.passwordReader == nil {
		c.passwordReader = new(stdInPasswordReader)
	}

	c.initGlobalFlags()

	if err := c.initProjectCmd(); err != nil {
		return nil, err
	}
	if err := c.initReleaseCmd(); err != nil {
		return nil, err
	}
	if err := c.initProviderCmd(); err != nil {
		return nil, err
	}

	if err := c.initSlackCmd(); err != nil {
		return nil, err
	}
	if err := c.initTelegramCmd(); err != nil {
		return nil, err
	}
	if err := c.initDiscordCmd(); err != nil {
		return nil, err
	}
	if err := c.initHangoutsChatCmd(); err != nil {
		return nil, err
	}
	if err := c.initMicrosoftTeamsCmd(); err != nil {
		return nil, err
	}
	if err := c.initWebhookCmd(); err != nil {
		return nil, err
	}

	c.initConfigureCmd()
	if err := c.initGetAuthKeyCmd(); err != nil {
		return nil, err
	}
	if err := c.initAuthCmd(); err != nil {
		return nil, err
	}
	c.initVersionCmd()
	return c, nil
}

func (c *command) Execute() (err error) {
	return c.root.Execute()
}

// Execute parses command line arguments and runs appropriate functions.
func Execute() (err error) {
	c, err := newCommand()
	if err != nil {
		return err
	}
	return c.Execute()
}

func (c *command) initGlobalFlags() {
	globalFlags := c.root.PersistentFlags()
	globalFlags.StringVar(&c.cfgFile, "config", "", "config file (default is $HOME/.newreleases.yaml)")
}

func (c *command) initConfig(requireConfigFileIfSet bool) (err error) {
	config := viper.New()
	configName := ".newreleases"
	if c.cfgFile != "" {
		// Use config file from the flag.
		config.SetConfigFile(c.cfgFile)
	} else {
		// Find home directory.
		if err := c.setHomeDir(); err != nil {
			return err
		}
		// Search config in home directory with name ".newreleases" (without extension).
		config.AddConfigPath(c.homeDir)
		config.SetConfigName(configName)
		requireConfigFileIfSet = false
	}

	// Environment
	config.SetEnvPrefix("newreleases")
	config.AutomaticEnv() // read in environment variables that match
	config.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	if c.homeDir != "" && c.cfgFile == "" {
		c.cfgFile = filepath.Join(c.homeDir, configName+".yaml")
	}

	// If a config file is found, read it in.
	if err := config.ReadInConfig(); err != nil {
		if requireConfigFileIfSet {
			return err
		}
		var e viper.ConfigFileNotFoundError
		if !errors.As(err, &e) && !os.IsNotExist(err) {
			return err
		}
	}
	c.config = config
	return nil
}

func (c *command) setHomeDir() (err error) {
	if c.homeDir != "" {
		return
	}
	dir, err := homedir.Dir()
	if err != nil {
		return err
	}
	c.homeDir = dir
	return nil
}

func (c *command) writeConfig(cmd *cobra.Command, authKey string) (err error) {
	c.config.Set(optionNameAuthKey, strings.TrimSpace(authKey))
	err = c.config.WriteConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		err = c.config.SafeWriteConfigAs(c.cfgFile)
	}
	return err
}

func newTable(w io.Writer) (table *tablewriter.Table) {
	table = tablewriter.NewWriter(w)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("   ")
	table.SetNoWhiteSpace(true)
	return table
}

func yesNo(b bool) (s string) {
	if b {
		return "yes"
	}
	return "no"
}
