// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

var RootCmd = rootCmd

func SetTestHomeDir(dir string) {
	testHomeDir = dir
}

func SetCfgFile(filename string) (reset func()) {
	reset = NewResetCfgFileFunc()
	cfgFile = filename
	return reset
}

func NewResetCfgFileFunc() (reset func()) {
	orig := cfgFile
	return func() { cfgFile = orig }
}

type PasswordReader = passwordReader

func SetCMDPasswordReader(new PasswordReader) (orig PasswordReader) {
	orig = cmdPasswordReader
	cmdPasswordReader = new
	return orig
}

type AuthKeysGetter = authKeysGetter

func SetCMDAuthKeysGetter(new AuthKeysGetter) (orig AuthKeysGetter) {
	orig = cmdAuthKeysGetter
	cmdAuthKeysGetter = new
	return orig
}

type AuthService = authService

func SetCMDAuthService(new AuthService) (orig AuthService) {
	orig = cmdAuthService
	cmdAuthService = new
	return orig
}
