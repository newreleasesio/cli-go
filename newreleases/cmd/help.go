// Copyright (c) 2020, NewReleases CLI AUTHORS.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file. {

package cmd

var configurationHelp = `Initial configuration:

  This tool needs to authenticate to NewReleases API using a secret Auth Key
  that can be generated on the service settings web pages.

  The key can be stored permanently by issuing interactive commands:

    newreleases configure

  or

    newreleases get-auth-key

  or it can be provided as the command line argument flag --auth-key on every
  newreleases tool execution.`
