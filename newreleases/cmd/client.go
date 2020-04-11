// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"context"
	"errors"
	"net/url"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"newreleases.io/newreleases"
)

func (c *command) getClient(cmd *cobra.Command) (client *newreleases.Client, err error) {
	if c.client != nil {
		return c.client, nil
	}

	authKey := c.config.GetString(optionNameAuthKey)
	if authKey == "" {
		cmd.Println(configurationHelp)
		cmd.Println()
		return nil, errors.New("auth key not configured")
	}
	o, err := newClientOptions(cmd)
	if err != nil {
		return nil, err
	}
	c.client = newreleases.NewClient(authKey, o)
	return c.client, nil
}

func newClientOptions(cmd *cobra.Command) (o *newreleases.ClientOptions, err error) {
	v, err := cmd.Flags().GetString(optionNameAPIEndpoint)
	if err != nil {
		return nil, err
	}
	var baseURL *url.URL
	if v != "" {
		baseURL, err = url.Parse(v)
		if err != nil {
			return nil, err
		}
	}
	return &newreleases.ClientOptions{BaseURL: baseURL}, nil
}

func newClientContext(config *viper.Viper) (ctx context.Context, cancel context.CancelFunc) {
	return context.WithTimeout(context.Background(), config.GetDuration(optionNameTimeout))
}

func addClientFlags(cmd *cobra.Command) (err error) {
	flags := cmd.Flags()
	flags.String(optionNameAuthKey, "", "API auth key")
	flags.Duration(optionNameTimeout, 30*time.Second, "API request timeout")
	flags.String(optionNameAPIEndpoint, "", "API Endpoint")
	if err := flags.MarkHidden(optionNameAPIEndpoint); err != nil {
		return err
	}
	return nil
}

func addClientConfigOptions(cmd *cobra.Command, config *viper.Viper) (err error) {
	flags := cmd.Flags()
	if err := config.BindPFlag(optionNameAuthKey, flags.Lookup(optionNameAuthKey)); err != nil {
		return err
	}
	if err := config.BindPFlag(optionNameTimeout, flags.Lookup(optionNameTimeout)); err != nil {
		return err
	}
	return nil
}
