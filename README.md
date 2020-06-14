![Latest Release Version](https://img.shields.io/github/v/release/cinarmert/doclogs)
[![Go Report Card](https://goreportcard.com/badge/github.com/cinarmert/doclogs)](https://goreportcard.com/report/github.com/cinarmert/doclogs)
![Go CI](https://github.com/cinarmert/doclogs/workflows/Go%20CI/badge.svg)

# `doclogs`

**`doclogs`** helps you view multiple docker container logs in the same splitted terminal session.

![doclogs demo gif](img/doclogs-demo.gif)

**`doclogs`** is a minimalistic cli tool. See the usage below!

```
Doclogs provides a user interface for the logs from multiple docker containers.

Usage:
  doclogs [OPTIONS] [CONTAINERS...] [flags]

Flags:
  -f, --follow    follow the stream of logs
  -h, --help      help for doclogs
  -v, --verbose   print debug logs
```

# Installation

**`doclogs`** is available on macOS, Linux and Windows. You can find the binaries in [**Releases &rarr;**](https://github.com/cinarmert/doclogs/releases).

## Brew

**`doclogs`** is already available in Homebrew for an easier installation.

```
brew tap cinarmert/doclogs
brew install doclogs
```
