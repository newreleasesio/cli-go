// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

type passwordReader interface {
	ReadPassword() (password string, err error)
}

type stdInPasswordReader struct{}

func (stdInPasswordReader) ReadPassword() (password string, err error) {
	v, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	return string(v), err
}

func terminalPromptPassword(cmd *cobra.Command, r passwordReader, title string) (password string, err error) {
	cmd.Print(title + ": ")
	password, err = r.ReadPassword()
	cmd.Println()
	if err != nil {
		return "", err
	}
	return password, nil
}

func terminalPrompt(cmd *cobra.Command, reader interface{ ReadString(byte) (string, error) }, title string) (value string, err error) {
	cmd.Print(title + ": ")
	value, err = reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(value), nil
}
