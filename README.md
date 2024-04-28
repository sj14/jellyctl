# jellyctl

A CLI for managing your Jellyfin server.

## Installation

### Precompiled Binaries

Binaries are available for all major platforms. See the [releases](https://github.com/sj14/jellyctl/releases) page.

### Homebrew

Using the [Homebrew](https://brew.sh/) package manager for macOS:

``` text
brew install sj14/tap/jellyctl
```

### Manually

It's also possible to install via `go install`:

```console
go install github.com/sj14/jellyctl@latest
```

## Usage

```
NAME:
   jellyctl - Manage Jellyfin from the CLI

USAGE:
   jellyctl [global options] command [command options] 

VERSION:
   undefined

COMMANDS:
   activity  List activities
   system    Manage system
   user      Manage users
   library   Manage media libraries
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --url value    URL of the Jellyfin server (default: "http://127.0.0.1:8096") [$JELLYCTL_URL]
   --token value  API token [$JELLYCTL_TOKEN]
   --help, -h     show help
   --version, -v  print the version
```

Generate an API token in the Jellyfin WebUI at Administration -> Overview -> Advanced -> API Token.
