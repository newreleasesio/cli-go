// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "newreleases",
	Short: "Release tracker for software engineers",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cmd.Help(); err != nil {
			handleError(cmd, err)
		}
	},
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
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home := findHomeDir()
		// Search config in home directory with name ".newreleases" (without extension).
		viper.AddConfigPath(home)
		configName := ".newreleases"
		viper.SetConfigName(configName)
		configType := "yaml"
		viper.SetConfigType(configType)
		cfgFile = filepath.Join(home, configName+"."+configType)
	}

	// Environment
	viper.SetEnvPrefix("newreleases")
	viper.AutomaticEnv() // read in environment variables that match
	viper.Set(optionNameAuthKey, "")

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		var e viper.ConfigFileNotFoundError
		if !errors.As(err, &e) {
			// Function cobra.OnInitialize does not provide error propagation, so
			// handle the error as it would be handled in the main package.
			fmt.Fprintln(os.Stderr, err)
			exit(1)
			return
		}
	}
}

func findHomeDir() (dir string) {
	if testHomeDir != "" {
		return testHomeDir
	}
	dir, err := homedir.Dir()
	if err != nil {
		// Function cobra.OnInitialize does not provide error propagation, so
		// handle the error as it would be handled in the main package.
		fmt.Fprintln(os.Stderr, err)
		exit(1)
		return
	}
	return dir
}

var testHomeDir = ""
