// Copyright (c) 2019, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmd

import "runtime/debug"

// automatically set on release
// and updated with vcs revision on init
var version = "0.0.0"

func init() {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return
	}

	var revision string
	var dirtyBuild bool
	for _, s := range info.Settings {
		switch s.Key {
		case "vcs.revision":
			revision = s.Value
		case "vcs.modified":
			dirtyBuild = s.Value == "true"
		}
	}

	if len(revision) == 0 {
		return
	}

	if len(revision) > 7 {
		revision = revision[:7]
	}

	version += "-" + revision
	if dirtyBuild {
		version += "-dirty"
	}
}

func Version() string {
	return version
}
