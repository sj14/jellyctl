# NAME

jellyctl - Manage Jellyfin from the CLI

# SYNOPSIS

jellyctl

```
[--help|-h]
[--token]=[value]
[--url]=[value]
[--version|-v]
```

**Usage**:

```
jellyctl [GLOBAL OPTIONS] [command [COMMAND OPTIONS]] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--help, -h**: show help

**--token**="": API token

**--url**="": URL of the Jellyfin server (default: http://127.0.0.1:8096)

**--version, -v**: print the version


# COMMANDS

## activity

List activities

**--after**="": only logs after the given time (default: 0001-01-01 00:00:00 +0000 UTC)

**--help, -h**: show help

**--json, -j**: print output as JSON

**--limit**="": limit the output logs (default: 15)

**--start**="": start at the given index (default: 0)

### help, h

Shows a list of commands or help for one command

## system

Manage the system

**--help, -h**: show help

### shutdown

Stop the Jellyfin process

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### restart

Restart the Jellyfin process

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### info

Shows system information

**--help, -h**: show help

**--public**: show public info which won't need a token

#### help, h

Shows a list of commands or help for one command

### backup

Export some data (EXPERIMENTAL)

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### restore

Import played and favourite information (based on the user name not user ID!) (EXPERIMENTAL)

**--help, -h**: show help

**--unfav**: unfavorite media when this is the backup state

**--unplay**: mark media as unplayed when this is the backup state

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## user

Manage users

**--help, -h**: show help

### create

Add a user

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### delete

Remove a user

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### enable

Enable a user

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### disable

Disable a user

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### set-admin

Promote admin privileges

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### unset-admin

Revoke admin privileges

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### set-hidden

Hide user from login page

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### unset-hidden

Expose user on login page

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### list

Shows users

**--help, -h**: show help

**--json, -j**: print output as JSON

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## library

Manage media libraries

**--help, -h**: show help

### scan

Trigger a rescan of all libraries

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### unscraped

List entries which were not scraped

**--help, -h**: show help

**--json, -j**: print output as JSON

**--types**="": filter media types (default: [Movie Series])

#### help, h

Shows a list of commands or help for one command

### search

Search throught the library

**--help, -h**: show help

**--json, -j**: print output as JSON

**--types**="": filter media types (default: [Movie Series])

#### help, h

Shows a list of commands or help for one command

### duplicates

List duplicates in the library

**--help, -h**: show help

**--json, -j**: print output as JSON

**--types**="": filter media types (default: [Movie Series])

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## key

Manage API keys

**--help, -h**: show help

### list

Show keys

**--help, -h**: show help

**--json, -j**: print output as JSON

#### help, h

Shows a list of commands or help for one command

### create

Add a new key

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### delete

Remove a key

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## task

Manage schedule tasks

**--help, -h**: show help

### list

Show tasks

**--help, -h**: show help

**--json, -j**: print output as JSON

#### help, h

Shows a list of commands or help for one command

### start

Start task

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### stop

Stop task

**--help, -h**: show help

#### help, h

Shows a list of commands or help for one command

### help, h

Shows a list of commands or help for one command

## help, h

Shows a list of commands or help for one command
