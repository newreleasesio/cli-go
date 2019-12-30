// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import "io"

type (
	Command         = command
	Option          = option
	PasswordReader  = passwordReader
	AuthService     = authService
	AuthKeysGetter  = authKeysGetter
	ProviderService = providerService
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

func WithProviderService(s ProviderService) func(c *Command) {
	return func(c *Command) {
		c.providerService = s
	}
}
