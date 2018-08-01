# chord
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=for-the-badge)](https://godoc.org/github.com/azillion/chord)
[![Github All Releases](https://img.shields.io/github/downloads/azillion/chord/total.svg?style=for-the-badge)](https://github.com/azillion/chord/releases)

A Discord TUI for direct messaging.

 * [Installation](README.md#installation)
      * [Binaries](README.md#binaries)
      * [Via Go](README.md#via-go)
 * [Usage](README.md#usage)
   * [Auth](README.md#auth)
 * [Uses](README.md#uses)

## Installation

#### Binaries

For installation instructions from binaries please visit the [Releases Page](https://github.com/azillion/chord/releases).

#### Via Go

```console
$ go get github.com/azillion/chord
```

## Usage

```console
$ chord -h
chord -  A Discord TUI for direct messaging.

Usage: chord <command>

Flags:

  -d              enable debug logging (default: false)
  -e, --email     email for Discord account (default: <none>)
  -p, --password  password for Discord account (default: <none>)
  -t, --token     token for Discord account (default: <none>)

Commands:

  config   Configure chord Discord settings.
  ls       List available Discord channels.
  tui      TUI for sending and receiving Discord DMs
  version  Show the version information.
```

**NOTE:** Be aware that you may need to login to [Discord](https://discordapp.com/) on a browser before you are able to run `chord`

### Auth

`chord` will automatically try to parse your config credentials stored at `$HOME/.chord.config`, but if
not present, you can pass through flags directly or call `chord config`.

## Uses:
- https://github.com/genuinetools/pkg

- https://github.com/bwmarrin/discordgo

- https://github.com/marcusolsson/tui-go


Makefile from [genuinetools](https://github.com/genuinetools/reg/blob/master/Makefile)
