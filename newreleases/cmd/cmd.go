// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"newreleases.io/newreleases"
)

const (
	optionNameAuthKey     = "auth-key"
	optionNameTimeout     = "timeout"
	optionNameAPIEndpoint = "api-endpoint"
)

var cmdPasswordReader passwordReader = new(stdInPasswordReader)

type passwordReader interface {
	ReadPassword() (password string, err error)
}

type stdInPasswordReader struct{}

func (stdInPasswordReader) ReadPassword() (password string, err error) {
	v, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	return string(v), err
}

func terminalPrompt(cmd *cobra.Command, reader interface{ ReadString(byte) (string, error) }, title string) (value string, err error) {
	cmd.Print(title + ": ")
	value, err = reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(value), nil
}

func terminalPromptPassword(cmd *cobra.Command, title string) (password string, err error) {
	cmd.Print(title + ": ")
	password, err = cmdPasswordReader.ReadPassword()
	cmd.Println()
	if err != nil {
		return "", err
	}
	return password, nil
}

func writeConfig(cmd *cobra.Command, authKey string) (err error) {
	viper.Set(optionNameAuthKey, strings.TrimSpace(authKey))
	err = viper.WriteConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		err = viper.SafeWriteConfigAs(cfgFile)
	}
	return err
}

func newClient() (client *newreleases.Client, err error) {
	authKey := viper.GetString(optionNameAuthKey)
	if authKey == "" {
		return nil, errors.New("auth key not configured")
	}
	return newreleases.NewClient(authKey, newClientOptions()), nil
}

func addClientFlags(cmd *cobra.Command) {
	flags := cmd.Flags()
	flags.String(optionNameAuthKey, "", "API auth key")
	flags.Duration(optionNameTimeout, 30*time.Second, "API request timeout")
	flags.String(optionNameAPIEndpoint, "", "API Endpoint")
	must(flags.MarkHidden(optionNameAPIEndpoint))

	cobra.OnInitialize(func() {
		must(viper.BindPFlag(optionNameAuthKey, flags.Lookup(optionNameAuthKey)))
		must(viper.BindPFlag(optionNameTimeout, flags.Lookup(optionNameTimeout)))
	})
}

func newClientOptions() (o *newreleases.ClientOptions) {
	return &newreleases.ClientOptions{
		BaseURL: mustURLParse(viper.GetString(optionNameAPIEndpoint)),
	}
}

func newClientContext() (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), viper.GetDuration(optionNameTimeout))
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func mustURLParse(s string) (u *url.URL) {
	if s == "" {
		return nil
	}
	u, err := url.Parse(s)
	must(err)
	return u
}
