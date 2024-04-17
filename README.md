# jellyctl

A CLI for managing your Jellyfin server.

## Installation

Currently, installation is only possible via `go install`:

```console
go install github.com/sj14/jellyctl@latest
```

## Usage

```
NAME:
   jellyctl - Manage Jellyfin from the CLI

USAGE:
   jellyctl [global options] command [command options] 

COMMANDS:
   activity  List activities
   system    Show system information
   user      Manage users
   help, h   Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --url value    URL of the Jellyfin server (default: "http://127.0.0.1:8096") [$JELLYCTL_URL]
   --token value  API token [$JELLYCTL_TOKEN]
   --help, -h     show help
```
