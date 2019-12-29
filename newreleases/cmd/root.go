// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"path/filepath"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:           "newreleases",
	Short:         "Release tracker for software engineers",
	SilenceErrors: true,
	SilenceUsage:  true,
}

var cfgFile string

func init() {
	cobra.OnInitialize(initConfig)

	// Persistent flags, which are global for application.
	globalFlags := rootCmd.PersistentFlags()
	globalFlags.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.newreleases.yaml)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	configName := ".newreleases"
	var home string
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home = findHomeDir()
		// Search config in home directory with name ".newreleases" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(configName)
	}

	// Environment
	viper.SetEnvPrefix("newreleases")
	viper.AutomaticEnv() // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		var e viper.ConfigFileNotFoundError
		if !errors.As(err, &e) {
			// Function cobra.OnInitialize does not provide error propagation, so
			// handle the error as it would be handled in the main package.
			must(err)
		} else if home != "" {
			cfgFile = filepath.Join(home, configName+".yaml")
		}
	}
}

func findHomeDir() (dir string) {
	if testHomeDir != "" {
		return testHomeDir
	}
	dir, err := homedir.Dir()
	// Function cobra.OnInitialize does not provide error propagation, so
	// handle the error as it would be handled in the main package.
	must(err)
	return dir
}

// testHomeDir is set on test runs in order not to interfere with potential
// configuration in user dir of the user that runs tests.
var testHomeDir string
