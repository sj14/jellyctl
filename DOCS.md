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
jellyctl [GLOBAL OPTIONS] command [COMMAND OPTIONS] [ARGUMENTS...]
```

# GLOBAL OPTIONS

**--help, -h**: show help

**--token**="": API token (default: 82073e8e9239430fb728c0e40131f03b)

**--url**="": URL of the Jellyfin server (default: "http://127.0.0.1:8096")

**--version, -v**: print the version


# COMMANDS

## activity

List activities

**--after**="": only logs after the given time (default: 0001-01-01 00:00:00 +0000 UTC)

**--limit**="": limit the output logs (default: 15)

**--start**="": start at the given index (default: 0)

## system

Manage the system

### shutdown

Stop the Jellyfin process

### restart

Restart the Jellyfin process

### info

Shows system information

**--public**: show public info which won't need a token

## user

Manage users

### create

Add a user

### delete

Remove a user

### enable

Enable a user

### disable

Disable a user

### set-admin

Promote admin privileges

### unset-admin

Revoke admin privileges

### set-hidden

Hide user from login page

### unset-hidden

Expose user on login page

### list

Shows users

## library

Manage media libraries

### scan

Trigger a rescan of all libraries

### unscraped

List entries which were not scraped

**--episodes**: show unscraped episodes

**--movies**: show unscraped movies

**--seasons**: show unscraped seasons

**--shows**: show unscraped shows

### search

Search throught the library

**--episodes**: show unscraped episodes

**--movies**: show unscraped movies

**--seasons**: show unscraped seasons

**--shows**: show unscraped shows

## key

Manage API keys

### list

Show keys

### create

Add a new key

### delete

Remove a key

## help, h

Shows a list of commands or help for one command
