<a href="https://newreleases.io"><img src="https://newreleases.io/logo.svg" alt="NewReleases" width="576"></a>

A command line client for managing [NewReleases](https://newreleases.io) projects.

# Installation

NewReleases client binaries have no external dependencies and can be just copied and executed locally.

Binary downloads of the NewReleases client can be found on the [Releases page](https://github.com/newreleasesio/cli-go/releases/latest).

To install on macOS:

```sh
wget https://github.com/newreleasesio/cli-go/releases/latest/download/newreleases-darwin-amd64 -O /usr/local/bin/newreleases
chmod +x /usr/local/bin/newreleases
```

You may need additional privileges to write to `/usr/local/bin`, but the file can be saved at any location that you want.

Supported operating systems and architectures:

- macOS 64bit `darwin-amd64`
- Linux 64bit `linux-amd64`
- Linux 32bit `linux-386`
- Linux ARM 64bit `linux-arm64`
- Linux ARM 32bit `linux-armv6`
- Windows 64bit `windows-amd64`
- Windows 32bit `windows-386`

Deb and RPM packages are also built.

This tool is implemented using the Go programming language and can be also installed by issuing a `go get` command:

```sh
go get -u newreleases.io/cmd/newreleases
```

# Configuration

This tool needs to authenticate to NewReleases API using a secret Auth Key
that can be generated on the service settings web pages.

The key can be stored permanently by issuing interactive commands:

```sh
newreleases configure
```

or

```sh
newreleases get-auth-key
```

or it can be provided as the command line argument flag `--auth-key` on every newreleases command execution.

# Usage

## Getting help

NewReleases client and its commands have help pages associated with them that can be printed out with `-h` flag:

```sh
newreleases -h
newreleases get-auth-key -h
newreleases project add -h
```

## Working with projects

The base command for getting releases is `project` and it shows available sub-commands which are `list`, `search`, `get`, `add`, `update` and `remove`.

### List projects

Listing all added projects is paginated and a page can be specified with `--page` (short `-p`) flag:

```sh
newreleases project list
newreleases project list -p 2
```

Project can be filtered by provider:

```sh
newreleases project list --provider github
```

and the order can be specified with `--order` flag which can have values `updated`, `added` or `name`:

```sh
newreleases project list --order name
```

Projects can be searched by name with:

```sh
newreleases project search go
```

where `go` is the example of a search string.

### Search projects

Search results can be filtered by provider, just as listing can be with `--provider` flag:

```sh
newreleases project search go --provider github
```

### Get a project

Information about a specific project can be retrieved with:

```sh
newreleases project get github golang/go
```

or by a project id:

```sh
newreleases project get mdsbe60td5gwgzetyksdfeyxt4
```

### Add new project to track

A project can be added with:

```sh
newreleases project add github golang/go
```

But there is a number of options that can be set, as by default, none of the notifications are enabled.

To enable emailing:

```sh
newreleases project add github golang/go --email daily
```

Or to add Slack notifications as well, but exclude pre-releases:

```sh
newreleases project add github golang/go --email daily --slack td5gwxt4mdsbe6gzetyksdfey0 --exclude-prereleases
```

More details about options can be found on `add` sub-command help page:

```sh
newreleases project add -h
```

### Update project options

Updating a project options is also possible. It contains the same options as the `add` command with additional flags to remove some of them. More information about options can be found on `update` sub-command help page:

```sh
newreleases project update -h
```

It is important that only specified options will be changed. For example, specifying different Slack channels will not remove already set other options like Telegram or Email or exclusions.

```sh
newreleases project update github golang/go --slack td5gwxt4mdsbe6gzetyksdfey0
```

### Remove a project

To remove the project from tracking its releases:

```sh
newreleases project remove github golang/go
```

or by a project id:

```sh
newreleases project remove mdsbe60td5gwgzetyksdfeyxt4
```

## Getting releases

The base command for getting releases is `release` and it shows available sub-commands which are `list`, `get`, and `note`.

### List releases of a project

To list all releases in chronological order of one project:

```sh
newreleases release list github golang/go
```

where the first argument after `list` is the provider and the second one is the project name.

or by project id:

```sh
newreleases release list mdsbe60td5gwgzetyksdfeyxt4
```

where the only argument after `list` is the project ID.

Results are paginated and the requested page can be specified with `--page` (short `-p`) flag.

```sh
newreleases release list github golang/go -p 2
```

### Get a release information

To get information about only one release, there is the `get` sub-command:

```sh
newreleases release get github golang/go go1.13.5
```

```sh
newreleases release get mdsbe60td5gwgzetyksdfeyxt4 go1.13.5
```

The last argument is the version which is requested.

### Get a release note

To get a release note about a release, there is the `note` sub-command:

```sh
newreleases release note npm vue 2.6.11
```

```sh
newreleases release note gzetyksdfeyxt4mdsbe60td5gw 2.6.11
```

## Listing providers

NewReleases supports a number of clients and they can be listed with:

```sh
newreleases providers
```

To list only providers that you have project added from:

```sh
newreleases providers --added
```

This information can be useful when filtering projects by a provider.

## Listing available notification channels

Notification channels can be managed only over the service's web interface. With NewReleases CLI client, they can be listed to relate their IDs from the output from other commands with their names. Available commands:

```sh
newreleases slack
newreleases telegram
newreleases discord
newreleases hangouts-chat
newreleases microsoft-teams
newreleases webhook
```

# Versioning

To see the current version of the binary, execute:

```sh
newreleases version
```

Each version is tagged and the version is updated accordingly in `version.go` file.

# Contributing

We love pull requests! Please see the [contribution guidelines](CONTRIBUTING.md).

# License

This application is distributed under the BSD-style license found in the [LICENSE](LICENSE) file.
