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
   system    Manage the system
   user      Manage users
   library   Manage media libraries
   key       Manage API tokens
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --url value    URL of the Jellyfin server (default: "http://127.0.0.1:8096") [$JELLYCTL_URL]
   --token value  API token [$JELLYCTL_TOKEN]
   --help, -h     show help
   --version, -v  print the version
```

Generate an API token in the Jellyfin WebUI at Administration -> Overview -> Advanced -> API Token.

### System

```
NAME:
   jellyctl system - Manage the system

USAGE:
   jellyctl system command [command options] 

COMMANDS:
   shutdown  Stop the Jellyfin process
   restart   Restart the Jellyfin process
   info      Shows system information
   help, h   Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### User

```
NAME:
   jellyctl user - Manage users

USAGE:
   jellyctl user command [command options] 

COMMANDS:
   create        Add a user
   delete        Remove a user
   enable        Enable a user
   disable       Disable a user
   set-admin     Promote admin privileges
   unset-admin   Revoke admin privileges
   set-hidden    Hide user from login page
   unset-hidden  Expose user on login page
   list          Shows users
   help, h       Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### Library

```
NAME:
   jellyctl library - Manage media libraries

USAGE:
   jellyctl library command [command options] 

COMMANDS:
   scan       Trigger a rescan of all libraries
   unscraped  List entries which were not scraped
   help, h    Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```

### Key

```
NAME:
   jellyctl key - Manage API tokens

USAGE:
   jellyctl key command [command options] 

COMMANDS:
   list     Show keys
   create   Add a new key
   delete   Remove a key
   help, h  Shows a list of commands or help for one command

OPTIONS:
   --help, -h  show help
```
